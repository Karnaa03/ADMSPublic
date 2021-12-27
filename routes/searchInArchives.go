package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	ginoidc "git.solutions.im/Solutions.IM/ginOidc"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10/orm"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/model"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

type SearchRequest struct {
	BookletNumber    string `form:"BookletNumber"`
	DivisionNumber   string `form:"DivisionNumber"`
	DivisionName     string
	DistrictNumber   string `form:"DistrictNumber"`
	DistrictName     string
	UpazilaNumber    string `form:"UpazilaNumber"`
	UpazilaName      string
	UnionNumber      string `form:"UnionNumber"`
	UnionName        string
	MouzaNumber      string `form:"MouzaNumber"`
	MouzaName        string
	VillageNumber    string `form:"VillageNumber"`
	VillageName      string
	EANumber         string `form:"EANumber"`
	ENName           string
	RMONumber        string `form:"RMONumber"`
	RMOName          string
	BoxToCheckout    string `form:"checkout"`
	BoxToCheckin     string `form:"checkin"`
	ArchiveBoxNumber string `form:"ArchiveBoxNumber"`
}

func (s SearchRequest) isEmpty() bool {
	if s.BookletNumber == "" &&
		s.DivisionNumber == "" &&
		s.DistrictNumber == "" &&
		s.UpazilaNumber == "" &&
		s.UnionNumber == "" &&
		s.MouzaNumber == "" &&
		s.VillageNumber == "" &&
		s.EANumber == "" &&
		s.RMONumber == "" &&
		s.ArchiveBoxNumber == "" {
		return true
	}
	return false
}

func (srv *Server) searchInArchives(footer string) {
	srv.router.GET("/warehouse/searchInArchives.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		id, err := ginoidc.GetIdentity(c)
		if err != nil {
			log.Error(err)
			srv.searchInArchivesWithError(c, header, footer, "who are you ?", SearchRequest{})
			return
		}
		var searchRequest SearchRequest
		err = c.ShouldBind(&searchRequest)
		if err != nil {
			log.Error(err)
			srv.searchInArchivesWithError(c, header, footer, err.Error(), searchRequest)
		}

		if searchRequest.BoxToCheckout != "" {
			log.Debugf("checkout archive box : %s", searchRequest.BoxToCheckout)
			searchRequest.ArchiveBoxNumber = searchRequest.BoxToCheckout
			a, err := srv.Db.GetArchiveBox(searchRequest.BoxToCheckout)
			if err != nil {
				log.Error(err)
				srv.searchInArchivesWithError(c, header, footer, "unable to find this archive box for checkout", searchRequest)
			}
			err = a.CheckOut(id, &srv.Db)
			if err != nil {
				log.Error(err)
				srv.searchInArchivesWithError(c, header, footer, fmt.Sprintf("error when trying to checkout = %s", err), searchRequest)
			}
		}
		if searchRequest.BoxToCheckin != "" {
			log.Debugf("checkin archive box : %s", searchRequest.BoxToCheckout)
			searchRequest.ArchiveBoxNumber = searchRequest.BoxToCheckin
			a, err := srv.Db.GetArchiveBox(searchRequest.BoxToCheckin)
			if err != nil {
				log.Error(err)
				srv.searchInArchivesWithError(c, header, footer, "unable to find this archive box for checkin", searchRequest)
			}
			err = a.CheckIn(id, &srv.Db)
			if err != nil {
				log.Error(err)
				srv.searchInArchivesWithError(c, header, footer, fmt.Sprintf("error when trying to checkin = %s", err), searchRequest)
			}
		}
		srv.searchInArchivesOkWithData(c, header, footer, searchRequest)
	})

	srv.router.POST("/warehouse/searchInArchives.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		var searchRequest SearchRequest
		err := c.ShouldBind(&searchRequest)
		if err != nil {
			srv.searchInArchivesWithError(c, header, footer, err.Error(), SearchRequest{})
		}
		srv.searchInArchivesOkWithData(c, header, footer, searchRequest)
	})
}

func (srv *Server) searchInArchivesOkWithData(c *gin.Context, header, footer string, sr SearchRequest) {
	name := ginoidc.GetValue(c, "name")

	csr := SearchRequest{}
	csr, tableData, err := srv.searchInArchivesWithQuery(sr)

	var content = gin.H{
		"Name":             name,
		"Header":           template.HTML(header),
		"Footer":           template.HTML(footer),
		"TableData":        template.HTML(tableData),
		"BookletNumber":    sr.BookletNumber,
		"DivisionNumber":   sr.DivisionNumber,
		"DivisionName":     csr.DivisionName,
		"DistrictNumber":   sr.DistrictNumber,
		"DistrictName":     csr.DistrictName,
		"UpazilaNumber":    sr.UpazilaNumber,
		"UpazilaName":      csr.UpazilaName,
		"UnionNumber":      sr.UnionNumber,
		"UnionName":        csr.UnionName,
		"MouzaNumber":      sr.MouzaNumber,
		"MouzaName":        csr.MouzaName,
		"VillageNumber":    sr.VillageNumber,
		"VillageName":      csr.VillageName,
		"EANumber":         sr.EANumber,
		"EAName":           csr.ENName,
		"RMONumber":        sr.RMONumber,
		"RMOName":          csr.RMOName,
		"ArchiveBoxNumber": sr.ArchiveBoxNumber,
	}
	if err != nil {
		content["Alert"] = template.HTML(err.Error())
	}
	c.HTML(http.StatusOK, "searchInArchives.html", content)
}

func (srv *Server) searchInArchivesWithQuery(sr SearchRequest) (completedSR SearchRequest, data string, err error) {
	var booklets []model.Booklet
	completedSR = sr
	req := srv.Db.Conn.Model(&booklets).
		Relation("ArchiveBox").
		Relation("GeoCode")
	// reminder : geoCode is in format : district.upazilla.union.mouza.ca.rmo

	if sr.isEmpty() {
		return completedSR, "", nil
	}
	if sr.BookletNumber != "" {
		req.Where("Booklet.number = ?", sr.BookletNumber)
	}
	if sr.DivisionNumber != "" {
		v, err := strconv.Atoi(sr.DivisionNumber)
		if err != nil {
			return completedSR, data, err
		}
		req.Where("geo_code.division = ?", v)
	}
	if sr.DistrictNumber != "" {
		v, err := strconv.Atoi(sr.DistrictNumber)
		if err != nil {
			return completedSR, data, err
		}
		req.Where("geo_code.district = ?", v)
	}
	if sr.UpazilaNumber != "" {
		v, err := strconv.Atoi(sr.UpazilaNumber)
		if err != nil {
			return completedSR, data, err
		}
		req.Where("geo_code.upazilla = ?", v)
	}
	if sr.UnionNumber != "" {
		v, err := strconv.Atoi(sr.UnionNumber)
		if err != nil {
			return completedSR, data, err
		}
		req.Where("geo_code.union = ?", v)
	}
	if sr.MouzaNumber != "" {
		v, err := strconv.Atoi(sr.MouzaNumber)
		if err != nil {
			return completedSR, data, err
		}
		req.Where("geo_code.mouza = ?", v)
	}
	if sr.VillageNumber != "" {
		v, err := strconv.Atoi(sr.VillageNumber)
		if err != nil {
			return completedSR, data, err
		}
		req.Where("geo_code.village = ?", v)
	}
	if sr.EANumber != "" {
		v, err := strconv.Atoi(sr.EANumber)
		if err != nil {
			return completedSR, data, err
		}
		req.Where("geo_code.ca = ?", v)
	}
	if sr.RMONumber != "" {
		v, err := strconv.Atoi(sr.RMONumber)
		if err != nil {
			return completedSR, data, err
		}
		req.Where("geo_code.rmo = ?", v)
	}
	if sr.ArchiveBoxNumber != "" {
		req.Where("archive_box.number = ?", sr.ArchiveBoxNumber)
	}
	req.Where("Booklet.status = ?", "archived").
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			query = query.WhereOr("archive_box.status = ?", "archived").
				WhereOr("archive_box.status = ?", "checkedOut")
			return query, nil
		})

	err = req.Select()
	if err != nil {
		log.Errorf("error when trying to request db : %s", err)
		return
	}
	for _, booklet := range booklets {
		data += fmt.Sprintf(`
		<tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%d</td>
			<td>%d</td>
			<td>%d</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>
`,
			booklet.Number,
			booklet.GeoCodeID,
			booklet.ArchiveBoxNumber,
			booklet.GetArchiveBox().WarehouseRowNumber,
			booklet.GetArchiveBox().WarehouseShelfNumber,
			booklet.GetArchiveBox().WarehouseShelfLevelNumber,
			booklet.GetCheckedOutBy(),
			booklet.GetCheckedInBy(),
			booklet.GetActionLink())
	}
	if booklets != nil && len(booklets) >= 0 {
		if sr.DistrictNumber != "" {
			completedSR.DistrictName = booklets[0].GeoCode.NameDistrict
		}
		if sr.DivisionNumber != "" {
			completedSR.DivisionName = booklets[0].GeoCode.NameDivision
		}
		if sr.UpazilaNumber != "" {
			completedSR.UpazilaName = booklets[0].GeoCode.NameUpazilla
		}
		if sr.UnionNumber != "" {
			completedSR.UnionName = booklets[0].GeoCode.NameUnion
		}
		if sr.MouzaNumber != "" {
			completedSR.MouzaName = booklets[0].GeoCode.NameMouza
		}
		if sr.VillageNumber != "" {
			vn, err := strconv.Atoi(sr.VillageNumber)
			if err != nil {
				log.Errorf("unable to parse village number %s", sr.VillageNumber)
			} else {
				completedSR.VillageName = lookupVillageName(*booklets[0].GeoCode.Villages, vn)
			}
		}
		if sr.EANumber != "" {
			completedSR.ENName = booklets[0].GeoCode.NameCountingArea
		}
		if sr.RMONumber != "" {
			completedSR.RMOName = booklets[0].GeoCode.NameRMO
		}
	}
	return
}

func (srv *Server) searchInArchivesWithError(c *gin.Context, header, footer, alertMsg string, sr SearchRequest) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")

	// tableData := srv.searchInArchivesWithQuery(sr)

	c.HTML(http.StatusOK, "searchInArchives.html", gin.H{
		"Name":             name,
		"Header":           template.HTML(header),
		"Footer":           template.HTML(footer),
		"Alert":            template.HTML(alert),
		"BookletNumber":    sr.BookletNumber,
		"DivisionNumber":   sr.DivisionNumber,
		"DistrictNumber":   sr.DistrictNumber,
		"UpazilaNumber":    sr.UpazilaNumber,
		"UnionNumber":      sr.UnionNumber,
		"MouzaNumber":      sr.MouzaNumber,
		"VillageNumber":    sr.VillageNumber,
		"EANumber":         sr.EANumber,
		"RMONumber":        sr.RMONumber,
		"ArchiveBoxNumber": sr.ArchiveBoxNumber,
	})
}
