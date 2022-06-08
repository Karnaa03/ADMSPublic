package model

type TableData struct {
	Count uint
	Rmo   uint
}

type Aggregates struct {
	Geocode   string `pg:",pk,notnull"`
	Rmo       uint   `pg:",pk,notnull"`
	HhSno     uint
	Sex       uint
	Sex2      uint
	Sex3      uint
	Edu       uint
	Edu1      uint
	Edu2      uint
	Edu3      uint
	Edu4      uint
	Edu5      uint
	Edu6      uint
	Edu7      uint
	Edu8      uint
	Edu9      uint
	Edu10     uint
	Edu12     uint
	Edu15     uint
	Edu18     uint
	Occ       uint
	Occ2      uint
	Occ3      uint
	Occ4      uint
	Occ5      uint
	HhF       uint
	HhF2      uint
	HhA       uint
	HhA2      uint
	C01M      uint
	C01F      uint
	C01H      uint
	C02M      uint
	C02F      uint
	C02H      uint
	C03M      uint
	C03F      uint
	C03H      uint
	C04       float64
	C05       float64
	C06       float64
	C07       float64
	C08       float64
	C09       float64
	C10       float64
	C11       float64
	C12       float64
	C13       float64
	C14       float64
	C15       float64
	C16       float64
	C17       float64
	C18       float64
	Sf        uint
	Mf        uint
	Lf        uint
	C19       float64
	C20       float64
	C21       float64
	C22       float64
	C23       float64
	C24       float64
	C25       float64
	C26       float64
	C27       float64
	C28H      uint
	C28F      uint
	C29H      uint
	C29F      uint
	C30H      uint
	C30F      uint
	C31H      uint
	C31F      uint
	C32H      uint
	C32F      uint
	C33H      uint
	C33F      uint
	C34H      uint
	C34F      uint
	C35H      uint
	C35F      uint
	C36H      uint
	C36F      uint
	C37H      uint
	C37F      uint
	C38H      uint
	C38F      uint
	C39       uint
	C40       uint
	C41A      uint
	C41B      uint
	C42A      uint
	C42B      uint
	C43A      uint
	C43B      uint
	C44A      uint
	C44B      uint
	C45A      uint
	C45B      uint
	C45C      uint
	C46A      uint
	C46B      uint
	C47A      uint
	C47B      uint
	C48       uint
	C49       uint
	T101      float64
	T102      float64
	T103      float64
	T104      float64
	T105      float64
	T112      float64
	T113      float64
	T114      float64
	T121      float64
	T122      float64
	T123      float64
	T124      float64
	T125      float64
	T127      float64
	T128      float64
	T129      float64
	T130      float64
	T131      float64
	T132      float64
	T134      float64
	T135      float64
	T157      float64
	T158      float64
	T159      float64
	T160      float64
	T161      float64
	T167      float64
	T169      float64
	T175      float64
	T176      float64
	T177      float64
	T179      float64
	T182      float64
	T185      float64
	T106      float64
	T107      float64
	T108      float64
	T109      float64
	T110      float64
	T111      float64
	T115      float64
	T116      float64
	T117      float64
	T118      float64
	T119      float64
	T120      float64
	T126      float64
	T133      float64
	T136      float64
	T137      float64
	T138      float64
	T139      float64
	T140      float64
	T141      float64
	T142      float64
	T143      float64
	T144      float64
	T145      float64
	T146      float64
	T147      float64
	T148      float64
	T149      float64
	T150      float64
	T151      float64
	T152      float64
	T153      float64
	T154      float64
	T155      float64
	T156      float64
	T162      float64
	T163      float64
	T164      float64
	T165      float64
	T166      float64
	T168      float64
	T170      float64
	T171      float64
	T172      float64
	T173      float64
	T174      float64
	T178      float64
	T180      float64
	T181      float64
	T183      float64
	T184      float64
	T186      float64
	T187      float64
	T188      float64
	T189      float64
	T190      float64
	T191      float64
	T192      float64
	T193      float64
	T194      float64
	T195      float64
	T196      float64
	T197      float64
	T198      float64
	T199      float64
	T200      float64
	T201      float64
	T202      float64
	T203      float64
	P501A     float64
	P501B     uint
	P502A     float64
	P502B     uint
	P503A     float64
	P503B     uint
	P504A     float64
	P504B     uint
	P505A     float64
	P505B     uint
	P506A     float64
	P506B     uint
	P507A     float64
	P507B     uint
	P508A     float64
	P508B     uint
	P510A     float64
	P510B     uint
	P511A     float64
	P511B     uint
	P512A     float64
	P512B     uint
	P521A     float64
	P521B     uint
	P522A     float64
	P522B     uint
	P523A     float64
	P523B     uint
	P524A     float64
	P524B     uint
	P538A     float64
	P538B     uint
	P539A     float64
	P539B     uint
	P546A     float64
	P546B     uint
	P548A     float64
	P548B     uint
	P549A     float64
	P549B     uint
	P550A     float64
	P550B     uint
	P551A     float64
	P551B     uint
	P572A     float64
	P572B     uint
	P509A     float64
	P509B     uint
	P513A     float64
	P513B     uint
	P514A     float64
	P514B     uint
	P515A     float64
	P515B     uint
	P516A     float64
	P516B     uint
	P517A     float64
	P517B     uint
	P518A     float64
	P518B     uint
	P519A     float64
	P519B     uint
	P520A     float64
	P520B     uint
	P525A     float64
	P525B     uint
	P526A     float64
	P526B     uint
	P527A     float64
	P527B     uint
	P528A     float64
	P528B     uint
	P529A     float64
	P529B     uint
	P530A     float64
	P530B     uint
	P531A     float64
	P531B     uint
	P532A     float64
	P532B     uint
	P533A     float64
	P533B     uint
	P534A     float64
	P534B     uint
	P535A     float64
	P535B     uint
	P536A     float64
	P536B     uint
	P537A     float64
	P537B     uint
	P540A     float64
	P540B     uint
	P541A     float64
	P541B     uint
	P542A     float64
	P542B     uint
	P543A     float64
	P543B     uint
	P544A     float64
	P544B     uint
	P545A     float64
	P545B     uint
	P547A     float64
	P547B     uint
	P552A     float64
	P552B     uint
	P553A     float64
	P553B     uint
	P554A     float64
	P554B     uint
	P555A     float64
	P555B     uint
	P556A     float64
	P556B     uint
	P557A     float64
	P557B     uint
	P559A     float64
	P559B     uint
	P560A     float64
	P560B     uint
	P561A     float64
	P561B     uint
	P562A     float64
	P562B     uint
	P563A     float64
	P563B     uint
	P564A     float64
	P564B     uint
	P565A     float64
	P565B     uint
	P566A     float64
	P566B     uint
	P567A     float64
	P567B     uint
	P568A     float64
	P568B     uint
	P569A     float64
	P569B     uint
	P570A     float64
	P570B     uint
	P571A     float64
	P571B     uint
	P573A     float64
	P573B     uint
	P574A     float64
	P574B     uint
	P575A     float64
	P575B     uint
	P577A     float64
	P577B     uint
	P579A     float64
	P579B     uint
	P580A     float64
	P580B     uint
	P581A     float64
	P581B     uint
	P582A     float64
	P582B     uint
	P584A     float64
	P584B     uint
	P585A     float64
	P585B     uint
	C04Gtrhh  uint
	C05Gtrhh  uint
	C06Gtrhh  uint
	C07Gtrhh  uint
	C08Gtrhh  uint
	C11Gtrhh  uint
	C12Gtrhh  uint
	C13Gtrhh  uint
	C14Gtrhh  uint
	C16Gtrhh  uint
	C17Gtrhh  uint
	C18Gtrhh  uint
	C19Gtrhh  uint
	C20Gtrhh  uint
	C21Gtrhh  uint
	C22Gtrhh  uint
	C23Gtrhh  uint
	C24Gtrhh  uint
	C25Gtrhh  uint
	C26Gtrhh  uint
	C27Gtrhh  uint
	C02Mfarm  uint
	C02Ffarm  uint
	C02Hfarm  uint
	C03Mfarm  uint
	C03Ffarm  uint
	C03Hfarm  uint
	C07Farm   float64
	C08Farm   float64
	C19Farm   float64
	C04Smlf   float64
	C05Smlf   float64
	C06Smlf   float64
	C07Smlf   float64
	C08Smlf   float64
	C09Smlf   float64
	C10Smlf   float64
	C11Smlf   float64
	C12Smlf   float64
	C13Smlf   float64
	C14Smlf   float64
	C15Smlf   float64
	C16Smlf   float64
	C17Smlf   float64
	C18Smlf   float64
	C19Smlf   float64
	C20Smlf   float64
	C21Smlf   float64
	C22Smlf   float64
	C23Smlf   float64
	C24Smlf   float64
	C25Smlf   float64
	C26Smlf   float64
	C27Smlf   float64
	C28Gtrhh  uint
	C29Gtrhh  uint
	C30Gtrhh  uint
	C31Gtrhh  uint
	C32Gtrhh  uint
	C33Gtrhh  uint
	C34Gtrhh  uint
	C35Gtrhh  uint
	C36Gtrhh  uint
	C37Gtrhh  uint
	C38Gtrhh  uint
	C39Gtrhh  uint
	C40Gtrhh  uint
	C41Gtrhh  uint
	C42Gtrhh  uint
	C43Gtrhh  uint
	C44Gtrhh  uint
	C45Gtrhh  uint
	C46Gtrhh  uint
	C47Gtrhh  uint
	C48Gtrhh  uint
	C49Gtrhh  uint
	C21Smlfhh uint
	C22Smlfhh uint
	C23Smlfhh uint
	C24Smlfhh uint
	C25Smlfhh uint
	C26Smlfhh uint
	C27Smlfhh uint
	C04Smlfhh uint
	C05Smlfhh uint
	C06Smlfhh uint
	C07Smlfhh uint
	C08Smlfhh uint
	C09Smlfhh uint
	C10Smlfhh uint
	C11Smlfhh uint
	C12Smlfhh uint
	C13Smlfhh uint
	C14Smlfhh uint
	C15Smlfhh uint
	C16Smlfhh uint
	C17Smlfhh uint
	C18Smlfhh uint
	C19Smlfhh uint
	C20Smlfhh uint
}
