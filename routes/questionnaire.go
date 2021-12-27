package routes

import (
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	ginoidc "git.solutions.im/Solutions.IM/ginOidc"
	"git.solutions.im/XeroxAgriCensus/AgriInject/goPg"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/oleiade/reflections"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/model"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

func (srv *Server) questionnaire() {
	srv.router.GET("/adms/questionnaire.html", func(c *gin.Context) {
		no := c.Query("QNumber")
		srv.questionnaireOkWithData(c, no, "")
	})

	srv.router.GET("/adms/getQuestionnaire", func(c *gin.Context) {
		tlno := c.Query("TlNumber")
		qno := c.Query("QNumber")
		content, err := srv.S3.Get(srv.Config.S3Config.Bucket, fmt.Sprintf("%s/%s.pdf", tlno, qno))
		if err != nil && err.Error() != io.EOF.Error() {
			log.Error(err)
		}
		c.Data(http.StatusOK, "application/pdf", content)
	})

	srv.router.GET("/adms/crop", func(context *gin.Context) {
		term := context.Query("query")
		s := struct {
			Query       string   `json:"query"`
			Suggestions []string `json:"suggestions"`
		}{}
		if term != "" {
			s.Query = term
			err := srv.Db.Conn.Model((*model.Crop)(nil)).
				ColumnExpr("distinct (id || ' - ' || name)").
				Where("(id || ' - ' || name like ?)", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%"))).
				Select(&s.Suggestions)
			if err != nil {
				log.Error(err)
			}
		}
		context.JSON(http.StatusOK, s)
	})

	srv.router.GET("/adms/tree", func(context *gin.Context) {
		term := context.Query("query")
		s := struct {
			Query       string   `json:"query"`
			Suggestions []string `json:"suggestions"`
		}{}
		if term != "" {
			s.Query = term
			err := srv.Db.Conn.Model((*model.Tree)(nil)).
				ColumnExpr("distinct (id || ' - ' || name)").
				Where("(id || ' - ' || name like ?)", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%"))).
				Select(&s.Suggestions)
			if err != nil {
				log.Error(err)
			}
		}
		context.JSON(http.StatusOK, s)
	})

	srv.router.POST("/adms/questionnaire.html", func(c *gin.Context) {
		q := goPg.Questionnaire{}
		no := c.Query("QNumber")
		err := c.Bind(&q)
		if err != nil {
			errMsg := fmt.Sprintf("error when trying to update questionnaire : %s", err)
			log.Error(errMsg)
			srv.questionnaireOkWithData(c, no, errMsg)
			return
		}
		// Issue 125 : Disable temporarily the village code check
		// geoCode, err := srv.Db.GetGeoCode(q.GeocodeID)
		// if err != nil {
		// 	srv.questionnaireOkWithData(c, no, fmt.Sprintf("Unable to find GeoCode with id : %s", q.GeocodeID))
		// 	return
		// }
		// if geoCode.Villages == nil || lookupVillageName(*geoCode.Villages, q.Village) == "" {
		// 	srv.questionnaireOkWithData(c, no, fmt.Sprintf("Village %d not found in reference data", q.Village))
		// 	return
		// }
		q.QuestionnaireNum = no
		srv.updateQuestionnaire(q, c)
	})
}

func (srv *Server) updateQuestionnaire(q goPg.Questionnaire, c *gin.Context) {
	user, err := ginoidc.GetIdentity(c)
	if err != nil {
		log.Error(err)
	}

	currentVersion := goPg.Questionnaire{}
	tx, err := srv.Db.Conn.Begin()
	if err != nil {
		log.Errorf("unable to open new transaction : %s", err)
	}
	currentVersion.QuestionnaireNum = q.QuestionnaireNum
	err = tx.Model(&currentVersion).WherePK().Select()
	if err != nil {
		log.Errorf("unable to get current version of the questionnaire : %s", err)
		_ = tx.Rollback()
		srv.questionnaireOkWithData(c, q.QuestionnaireNum, fmt.Sprintf("Error during questionnaire update, please contact your administrator\n%s", err))
		return
	}
	overWritten := goPg.OverwrittenQuestionnaire{
		UpdateDate:    time.Now(),
		UpdatedBy:     user.Id,
		Questionnaire: currentVersion,
	}
	_, err = tx.Model(&overWritten).Insert()
	if err != nil {
		log.Errorf("unable to insert current version in OverWritten table : %s", err)
		_ = tx.Rollback()
		srv.questionnaireOkWithData(c, q.QuestionnaireNum, fmt.Sprintf("Error during questionnaire update, please contact your administrator\n%s", err))
		return
	}

	// prepare new version
	q.QuestionnaireNum = currentVersion.QuestionnaireNum
	q.QuestionnaireEmpty = currentVersion.QuestionnaireEmpty
	q.FormName = currentVersion.FormName
	q.BookletNumber = currentVersion.BookletNumber
	q.TallySheetNo = currentVersion.TallySheetNo
	q.IsUnion = currentVersion.IsUnion
	NormalizeIdName(&q, table14Additional)
	NormalizeIdName(&q, table15Additional)
	res, err := tx.Model(&q).WherePK().Update()
	if err != nil {
		log.Errorf("unable to update questionnaire : %s", err)
		_ = tx.Rollback()
		srv.questionnaireOkWithData(c, q.QuestionnaireNum, fmt.Sprintf("Error during questionnaire update, please contact your administrator\n,%s", err))
		return
	}
	if res.RowsAffected() != 1 {
		log.Errorf("number of row affected : %d", res.RowsAffected())
		_ = tx.Rollback()
		srv.questionnaireOkWithData(c, q.QuestionnaireNum, fmt.Sprintf("Error during questionnaire update, please contact your administrator\n%s", err))
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Errorf("error during questionnaire update commit : %s", err)
		srv.questionnaireOkWithData(c, q.QuestionnaireNum, fmt.Sprintf("Error during questionnaire update, please contact your administrator\n%s", err))
		return
	}
	srv.questionnaireOkWithData(c, q.QuestionnaireNum, "")
}

func NormalizeIdName(q *goPg.Questionnaire, table [][]Field) {
	const errMsg = "error when trying to normalize ID / Name for additional surface : %s"
	for _, field := range table {
		v1, err := reflections.GetField(q, field[0].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}

		IdName := strings.Split(v1.(string), "-")
		if len(IdName) > 1 {
			err = reflections.SetField(q, field[0].name, strings.TrimSpace(IdName[1]))
			if err != nil {
				log.Errorf(errMsg, err)
			}
			id, err := strconv.Atoi(strings.TrimSpace(IdName[0]))
			if err != nil {
				log.Errorf(errMsg, err)
			}
			err = reflections.SetField(q, field[1].name, id)
			if err != nil {
				log.Errorf(errMsg, err)
			}
		}
	}
}

func checkTable9BusinessRules(q goPg.Questionnaire) (errs []error) {
	if q.TotalMemberMale < q.PeopleBetween10to14Male+q.PeopleAbove15Male {
		genre := "Male"
		errs = append(errs, fmt.Errorf("\"01 Total member of house %s\" must be less than \"02 People between 10 to 14yo %s\" + \"03 People above 15 %s\"", genre, genre, genre))
	}
	if q.TotalMemberFemale < q.PeopleBetween10to14Female+q.PeopleAbove15Female {
		genre := "Female"
		errs = append(errs, fmt.Errorf("\"01 Total member of house %s\" must be less than \"02 People between 10 to 14yo %s\" + \"03 People above 15 %s\"", genre, genre, genre))
	}
	if q.TotalMemberHijra < q.PeopleBetween10to14Hijra+q.PeopleAbove15Hijra {
		genre := "Hijra"
		errs = append(errs, fmt.Errorf("\"01 Total member of house %s\" must be less than \"02 People between 10 to 14yo %s\" + \"03 People above 15 %s\"", genre, genre, genre))
	}
	return
}

func checkTable10BusinessRules(q goPg.Questionnaire) (errs []error) {
	if q.TotalLand == 0 && q.CultivatedLand > 0.05 {
		if q.TotalLand != q.CultivatedLand {
			errs = append(errs, fmt.Errorf("04 Total own land must be equal to 18 Cultivated Land"))
		}
	}
	if q.LandGiven > q.TotalLand {
		errs = append(errs, fmt.Errorf("04 Total own land must be less or equal than 05 Land given"))
	}
	if !almostEqual(q.OperatingLand, q.TotalLand-q.LandGiven+q.LandTaken) {
		errs = append(errs, fmt.Errorf("07 Total operating Land must be equal to 04 Total own land - 05 Land given + 06 Land taken"))
	}
	if q.ResidenceLand > 30.00 {
		errs = append(errs, fmt.Errorf("08 Land for residence must be less than 30 acres"))
	}
	if q.BusinessEntityLand > 30.00 {
		errs = append(errs, fmt.Errorf("09 Land used for business entity must be less than 30 acres"))
	}
	if q.SinkChannelLand > 30.00 {
		errs = append(errs, fmt.Errorf("10 Land in sink-channel, bush etc must be less than 30 acres"))
	}
	if q.PermanentUnusedLand > 30.00 {
		errs = append(errs, fmt.Errorf("11 Permanent unused Land must be less than 30 acres"))
	}
	if !almostEqual(q.UncultivatedLand, q.ResidenceLand+q.BusinessEntityLand+q.SinkChannelLand+q.PermanentUnusedLand) {
		errs = append(errs, fmt.Errorf("12 Uncultivated Land must be equal to 08 Land for residence + 09 Land used for business entity + 10 Land in sink-channel, bush etc + 11 Permanent unused Land"))
	}
	if q.TempCultivatedLand > 30.00 {
		errs = append(errs, fmt.Errorf("13 Temporary cultivated Land must be less than 30 acres"))
	}
	if q.PondLand > 30.00 {
		errs = append(errs, fmt.Errorf("15 Land under pond/lake must be less than 30 acres"))
	}
	if q.NurseryLand > 30.00 {
		errs = append(errs, fmt.Errorf("16 Land under nursery must be less than 30 acres"))
	}
	if q.RecentUnusedLand > q.TempCultivatedLand {
		errs = append(errs, fmt.Errorf("17 Recent unused Land must be less or equal to 13 Temporary cultivated Land"))
	}
	if !almostEqual(q.CultivatedLand, q.TempCultivatedLand+q.PermanentCultivatedLand+q.PondLandNotBlank+q.NurseryLand+q.RecentUnusedLand) {
		errs = append(errs, fmt.Errorf("18 Cultivated Land must be equal to 13 Temporary cultivated Land + 14 Permanent cultivated land  + 15 Land under pond/lake + 16 Land under nursery + 17 Recent unused Land"))
	}
	if q.IrrigationLand > 30.00 {
		if q.IrrigationLand > q.TempCultivatedLand+q.PermanentCultivatedLand {
			errs = append(errs, fmt.Errorf("19 Land under irrigation must be less than 13 Temporary cultivated Land + 14 Permanent cultivated land"))
		}
	}
	if q.SaltCultivationLand > 60.00 {
		errs = append(errs, fmt.Errorf("20 Land under salt cultivation must be less than 60 acres"))
	}
	return
}

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold*(math.Abs(a)+math.Abs(b))
}

func (srv *Server) questionnaireOkWithData(c *gin.Context, no, alertMsg string) {
	user, err := ginoidc.GetIdentity(c)
	if err != nil {
		log.Error(err)
	}

	var alerts []string

	if alertMsg != "" {
		alert, err := templates.RenderAlert(alertMsg)
		if err != nil {
			log.Error(err)
		}
		alerts = append(alerts, alert)
	}

	q, err := srv.getQuestionnaire(no)
	if err != nil {
		log.Errorf("error when trying to get questionnaire data : %s", err)
	}

	errTable9 := checkTable9BusinessRules(q)
	if len(errTable9) != 0 {
		for _, err2 := range errTable9 {
			alert, err := templates.RenderWarning(err2.Error())
			if err != nil {
				log.Error(err)
			}
			alerts = append(alerts, alert)
		}
	}

	errTable10 := checkTable10BusinessRules(q)
	if len(errTable10) != 0 {
		for _, err2 := range errTable10 {
			alert, err := templates.RenderWarning(err2.Error())
			if err != nil {
				log.Error(err)
			}
			alerts = append(alerts, alert)
		}
	}

	var alertContent template.HTML
	for _, alert := range alerts {
		alertContent += template.HTML(alert)
	}

	c.HTML(http.StatusOK, "questionnaire.html", gin.H{
		"Alert":             alertContent,
		"User":              user,
		"Table9":            template.HTML(generateTableFormPeopleCount(table9, 9, q)),
		"Table10":           template.HTML(generateTableForm(table10, 10, q)),
		"Table11":           template.HTML(generateTableForm(table11, 11, q)),
		"Table12":           template.HTML(generateTableFormAnimalCount(table12, 12, q)),
		"Table13":           template.HTML(generateTableFormAnimalCount(table13, 13, q)),
		"Table14R":          template.HTML(generateTableForm(table14R, 14, q)),
		"Table14L":          template.HTML(generateTableForm(table14L, 14, q)),
		"Table14Additional": template.HTML(generateTableFormAdditionalCrop(table14Additional, 14, q, "crop")),
		"Table15":           template.HTML(generateTableFormSurfaceCount(table15, 15, q)),
		"Table15Additional": template.HTML(generateTableFormAdditionalSurfaceCount(table15Additional, 15, q, "tree")),
		"Table16":           template.HTML(generateTableFormMachine(table16, 16, q)),

		// Location
		"GeocodeID":    q.GeocodeID,
		"CountingArea": fmt.Sprintf("%03d - %s", q.GeoCode.CA, q.GeoCode.NameCountingArea),
		"District":     fmt.Sprintf("%02d - %s", q.GeoCode.District, q.GeoCode.NameDistrict),
		"Division":     fmt.Sprintf("%02d - %s", q.GeoCode.Division, q.GeoCode.NameDivision),
		"IsUnion":      q.IsUnion,
		"Mouza":        fmt.Sprintf("%03d - %s", q.GeoCode.Mouza, q.GeoCode.NameMouza),
		"RMO":          fmt.Sprintf("%d - %s", q.GeoCode.Rmo, q.GeoCode.NameRMO),
		"Union":        fmt.Sprintf("%03d - %s", q.GeoCode.Union, q.GeoCode.NameUnion),
		"Upazila":      fmt.Sprintf("%02d - %s", q.GeoCode.Upazilla, q.GeoCode.NameUpazilla),
		"Village":      fmt.Sprintf("%03d", q.Village),
		"VillageName":  lookupVillageName(*q.GeoCode.Villages, q.Village),

		// House hold information
		"TlNumber":         q.TallySheetNo,
		"QNumber":          q.QuestionnaireNum,
		"AgriLaborCode":    q.AgriLaborCode,
		"OwnerName":        q.OwnerName,
		"OwnerGender":      q.OwnerGender,
		"FatherName":       q.FatherName,
		"HouseSerial":      fmt.Sprintf("%03d", q.HouseSerial),
		"IsFishingRelated": q.IsFishingRelated,
		"MobileNo":         q.MobileNo,
		"OwnerEduCode":     fmt.Sprintf("%02d", q.OwnerEduCode),
		"OwnerProfesCode":  q.OwnerProfesCode,
	})
}

func lookupVillageName(villages model.Villages, id int) (names string) {
	if v, ok := villages.Villages[id]; ok {
		return fmt.Sprintf("%s %s", v.NameEn, v.NameBengali)
	} else {
		log.Warnf("village %d not found in %v", id, villages)
		return
	}
}

func (srv *Server) questionnaireWithError(c *gin.Context, alertMsg string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "questionnaire.html", gin.H{

		"Name":  name,
		"Alert": template.HTML(alert),
	})
}

func (srv *Server) getQuestionnaire(num string) (q goPg.Questionnaire, err error) {
	q.QuestionnaireNum = num
	err = srv.Db.Conn.Model(&q).Relation("GeoCode").WherePK().Select()
	return
}

func generateTableForm(table []Field, num int, dbQ goPg.Questionnaire) (form string) {
	form = fmt.Sprintf(`
<div class="x_title">
    <h2>Table %d (Surface in Acres)</h2>
    <div class="clearfix"></div>
</div>`, num)

	for _, field := range table {
		v, err := reflections.GetField(dbQ, field.name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field.name, err)
			continue
		}
		form += fmt.Sprintf(`
<div class="form-group">
    <label class="control-label col-md-6 col-sm-6 col-xs-12">
        %s
    </label>
    <div class="col-md-6 col-sm-6 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
</div>
`, field.label, field.inputMask, field.pattern, field.name, toFormString(v))
	}
	return
}

func generateTableFormPeopleCount(table [][]Field, num int, dbQ goPg.Questionnaire) (form string) {
	form = fmt.Sprintf(`
<div class="x_title">
    <h2>Table %d</h2>
    <div class="clearfix"></div>
</div>
	<div class="col-md-3 col-sm-3 col-xs-12"></div>
	    <label class="control-label col-md-3 col-sm-3 col-xs-12">Male</label>
		<label class="control-label col-md-3 col-sm-3 col-xs-12">Female</label>
		<label class="control-label col-md-3 col-sm-3 col-xs-12">Hijra</label>
	`, num)

	for _, field := range table {
		v1, err := reflections.GetField(dbQ, field[0].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v2, err := reflections.GetField(dbQ, field[1].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v3, err := reflections.GetField(dbQ, field[2].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		form += fmt.Sprintf(`
<div class="form-group">
    <label class="control-label col-md-3 col-sm-3 col-xs-12">
        %s
    </label>
    <div class="col-md-3 col-sm-3 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
	<div class="col-md-3 col-sm-3 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
	<div class="col-md-3 col-sm-3 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
</div>
`, field[0].label, field[0].inputMask, field[0].pattern, field[0].name, toFormString(v1),
			field[1].inputMask, field[1].pattern, field[1].name, toFormString(v2),
			field[2].inputMask, field[2].pattern, field[2].name, toFormString(v3),
		)

	}
	return
}

func generateTableFormMachine(table [][]Field, num int, dbQ goPg.Questionnaire) (form string) {
	form = fmt.Sprintf(`
<div class="x_title">
    <h2>Table %d</h2>
    <div class="clearfix"></div>
</div>
	<div class="col-md-3 col-sm-3 col-xs-12"></div>
		<label class="control-label col-md-3 col-sm-3 col-xs-12">Manual</label>
	    <label class="control-label col-md-3 col-sm-3 col-xs-12">Diesel</label>
		<label class="control-label col-md-3 col-sm-3 col-xs-12">Electric</label>
	`, num)

	for _, field := range table {
		v1, err := reflections.GetField(dbQ, field[0].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v2, err := reflections.GetField(dbQ, field[1].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v3, err := reflections.GetField(dbQ, field[2].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}

		var diesel, electric, manual string
		if field[0].value != "" {
			diesel = fmt.Sprintf(`<input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">`, field[0].inputMask, field[0].pattern, field[0].name, toFormString(v1))
		}
		if field[1].value != "" {
			electric = fmt.Sprintf(`<input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">`, field[1].inputMask, field[1].pattern, field[1].name, toFormString(v2))
		}
		if field[2].value != "" {
			manual = fmt.Sprintf(`<input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">`, field[2].inputMask, field[2].pattern, field[2].name, toFormString(v3))
		}

		form += fmt.Sprintf(`
<div class="form-group">
    <label class="control-label col-md-3 col-sm-3 col-xs-12">
        %s
    </label>
    <div class="col-md-3 col-sm-3 col-xs-12">
        %s
    </div>
	<div class="col-md-3 col-sm-3 col-xs-12">
        %s
    </div>
	<div class="col-md-3 col-sm-3 col-xs-12">
        %s
    </div>
</div>
`, field[0].label,
			diesel,
			electric,
			manual,
		)
	}
	return
}

func generateTableFormAdditionalCrop(table [][]Field, num int, dbQ goPg.Questionnaire, autoCompletePrefix string) (form string) {
	form = fmt.Sprintf(`
<div class="x_title">
    <h2>Table %d (Surface in Acres)</h2>
    <div class="clearfix"></div>
</div>
	<div class="col-md-3 col-sm-3 col-xs-12"></div>
	    <label class="control-label col-md-5 col-sm-5 col-xs-12">Id - Name</label>
		<label class="control-label col-md-4 col-sm-4 col-xs-12">Surface (Acre(s))</label>
	`, num)

	for i, field := range table {
		v1, err := reflections.GetField(dbQ, field[0].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v2, err := reflections.GetField(dbQ, field[1].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v3, err := reflections.GetField(dbQ, field[2].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}

		var addCropIdName string
		if v2 != 0 {
			addCropIdName = fmt.Sprintf("%d - %s", v2, toFormString(v1))
		}
		form += fmt.Sprintf(`
<div class="form-group">
    <label class="control-label col-md-3 col-sm-3 col-xs-12">
        %s
    </label>
    <div class="col-md-5 col-sm-5 col-xs-12">
        <input type="text"
               name="%s"
		       id="autocomplete-%s"	
               class="form-control"
               value="%s">
    </div>
	<div class="col-md-4 col-sm-4 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
</div>
`,
			field[0].label, field[0].name, fmt.Sprintf("%s-%d", autoCompletePrefix, i), addCropIdName,
			field[2].inputMask, field[2].pattern, field[2].name, toFormString(v3),
		)
	}
	return
}

func generateTableFormAdditionalSurfaceCount(table [][]Field, num int, dbQ goPg.Questionnaire, autoCompletePrefix string) (form string) {
	form = fmt.Sprintf(`
<div class="x_title">
    <h2>Table %d (Surface in Acres)</h2>
    <div class="clearfix"></div>
</div>
	<div class="col-md-3 col-sm-3 col-xs-12"></div>
	    <label class="control-label col-md-3 col-sm-3 col-xs-12">Id - Name</label>
		<label class="control-label col-md-3 col-sm-3 col-xs-12">Surface (Acre(s))</label>
		<label class="control-label col-md-3 col-sm-3 col-xs-12">Count</label>
	`, num)

	for i, field := range table {
		v1, err := reflections.GetField(dbQ, field[0].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v2, err := reflections.GetField(dbQ, field[1].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v3, err := reflections.GetField(dbQ, field[2].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v4, err := reflections.GetField(dbQ, field[3].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}

		var addCropIdName string
		if v2 != 0 {
			addCropIdName = fmt.Sprintf("%d - %s", v2, toFormString(v1))
		}
		form += fmt.Sprintf(`
<div class="form-group">
    <label class="control-label col-md-3 col-sm-3 col-xs-12">
        %s
    </label>
    <div class="col-md-3 col-sm-3 col-xs-12">
        <input type="text"
               name="%s"
		       id="autocomplete-%s"	
               class="form-control"
               value="%s">
    </div>
	<div class="col-md-3 col-sm-3 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
	<div class="col-md-3 col-sm-3 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
</div>
`,
			field[0].label, field[0].name, fmt.Sprintf("%s-%d", autoCompletePrefix, i), addCropIdName,
			field[3].inputMask, field[3].pattern, field[3].name, toFormString(v4),
			field[2].inputMask, field[2].pattern, field[2].name, toFormString(v3),
		)
	}
	return
}

func generateTableFormSurfaceCount(table [][]Field, num int, dbQ goPg.Questionnaire) (form string) {
	form = fmt.Sprintf(`
<div class="x_title">
    <h2>Table %d (Surface in Acres)</h2>
    <div class="clearfix"></div>
</div>
	<div class="col-md-4 col-sm-4 col-xs-12"></div>
	    <label class="control-label col-md-4 col-sm-4 col-xs-12">Surface</label>
		<label class="control-label col-md-4 col-sm-4 col-xs-12">Count</label>
	`, num)

	for _, field := range table {
		v1, err := reflections.GetField(dbQ, field[0].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v2, err := reflections.GetField(dbQ, field[1].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		form += fmt.Sprintf(`
<div class="form-group">
    <label class="control-label col-md-4 col-sm-4 col-xs-12">
        %s
    </label>
    <div class="col-md-4 col-sm-4 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
	<div class="col-md-4 col-sm-4 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
</div>
`, field[0].label, field[0].inputMask, field[0].pattern, field[0].name, toFormString(v1),
			field[1].inputMask, field[1].pattern, field[1].name, toFormString(v2),
		)

	}
	return
}

func generateTableFormAnimalCount(table [][]Field, num int, dbQ goPg.Questionnaire) (form string) {
	form = fmt.Sprintf(`
<div class="x_title">
    <h2>Table %d</h2>
    <div class="clearfix"></div>
</div>
	<div class="col-md-4 col-sm-4 col-xs-12"></div>
	    <label class="control-label col-md-4 col-sm-4 col-xs-12">At Home</label>
		<label class="control-label col-md-4 col-sm-4 col-xs-12">At Farm</label>
	`, num)

	for _, field := range table {
		v1, err := reflections.GetField(dbQ, field[0].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		v2, err := reflections.GetField(dbQ, field[1].name)
		if err != nil {
			log.Errorf("unable to get value for field %s : %s", field[0].name, err)
			continue
		}
		form += fmt.Sprintf(`
<div class="form-group">
    <label class="control-label col-md-4 col-sm-4 col-xs-12">
        %s
    </label>
    <div class="col-md-4 col-sm-4 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
	<div class="col-md-4 col-sm-4 col-xs-12">
        <input type="text"
               class="form-control"
               data-inputmask="%s"
               pattern="%s"
               name="%s"
               value="%s">
    </div>
</div>
`, field[0].label, field[0].inputMask, field[0].pattern, field[0].name, toFormString(v1),
			field[1].inputMask, field[1].pattern, field[1].name, toFormString(v2),
		)

	}
	return
}

func toFormString(v interface{}) string {
	switch t := v.(type) {
	case int:
		if t == 0 {
			return ""
		} else {
			return fmt.Sprintf("%05d", t)
		}
	case float64:
		if t == 0.0 {
			return ""
		} else {
			return fmt.Sprintf("%05.2f", t)
		}
	default:
		return fmt.Sprintf("%v", t)
	}
}
