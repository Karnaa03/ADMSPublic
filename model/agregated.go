package model

type Agregated struct {
	Geocode string `pg:",pk,notnull"`
	Rmo     uint   `pg:",pk,notnull"`
	HhSno   uint
	Sex     uint
	Sex2    uint
	Sex3    uint
	Edu     uint
	Edu1    uint
	Edu2    uint
	Edu3    uint
	Edu4    uint
	Edu5    uint
	Edu6    uint
	Edu7    uint
	Edu8    uint
	Edu9    uint
	Edu10   uint
	Edu12   uint
	Edu15   uint
	Edu18   uint
	Occ     uint
	Occ2    uint
	Occ3    uint
	Occ4    uint
	Occ5    uint
	HhF     uint
	HhF2    uint
	HhA     uint
	HhA2    uint
	C01M    uint
	C01F    uint
	C01H    uint
	C02M    uint
	C02F    uint
	C02H    uint
	C03M    uint
	C03F    uint
	C03H    uint
	C04     float64
	C05     float64
	C06     float64
	C07     float64
	C08     float64
	C09     float64
	C10     float64
	C11     float64
	C12     float64
	C13     float64
	C14     float64
	C15     float64
	C16     float64
	C17     float64
	C18     float64
	Sf      uint
	Mf      uint
	Lf      uint
	C19     float64
	C20     uint
	C21     float64
	C22     float64
	C23     float64
	C24     float64
	C25     uint
	C26     float64
	C27     float64
	C28H    uint
	C28F    uint
	C29H    uint
	C29F    uint
	C30H    uint
	C30F    uint
	C31H    uint
	C31F    uint
	C32H    uint
	C32F    uint
	C33H    uint
	C33F    uint
	C34H    uint
	C34F    uint
	C35H    uint
	C35F    uint
	C36H    uint
	C36F    uint
	C37H    uint
	C37F    uint
	C38H    uint
	C38F    uint
	C39     uint
	C40     uint
	C41A    uint
	C41B    uint
	C42A    uint
	C42B    uint
	C43A    uint
	C43B    uint
	C44A    uint
	C44B    uint
	C45A    uint
	C45B    uint
	C45C    uint
	C46A    uint
	C46B    uint
	C47A    uint
	C47B    uint
	C48     uint
	C49     uint
	T101    float64
	T102    float64
	T103    float64
	T104    float64
	T105    float64
	T112    float64
	T113    float64
	T114    float64
	T121    float64
	T122    float64
	T123    float64
	T124    float64
	T125    float64
	T127    float64
	T128    float64
	T129    float64
	T130    float64
	T131    float64
	T132    float64
	T134    float64
	T135    float64
	T157    float64
	T158    float64
	T159    float64
	T160    float64
	T161    float64
	T167    uint
	T169    float64
	T175    float64
	T176    float64
	T177    float64
	T179    float64
	T182    float64
	T185    float64
	T106    uint
	T107    uint
	T108    uint
	T109    uint
	T110    uint
	T111    uint
	T115    uint
	T116    uint
	T117    uint
	T118    uint
	T119    uint
	T120    uint
	T126    uint
	T133    uint
	T136    uint
	T137    uint
	T138    uint
	T139    uint
	T140    uint
	T141    uint
	T142    uint
	T143    uint
	T144    uint
	T145    uint
	T146    uint
	T147    uint
	T148    uint
	T149    uint
	T150    uint
	T151    uint
	T152    uint
	T153    uint
	T154    uint
	T155    uint
	T156    uint
	T162    uint
	T163    uint
	T164    uint
	T165    uint
	T166    uint
	T168    uint
	T170    uint
	T171    uint
	T172    uint
	T173    uint
	T174    uint
	T178    uint
	T180    uint
	T181    uint
	T183    uint
	T184    uint
	T186    uint
	T187    uint
	T188    uint
	T189    uint
	T190    uint
	T191    uint
	T192    uint
	T193    uint
	T194    uint
	T195    uint
	T196    uint
	T197    uint
	T198    uint
	T199    uint
	T200    uint
	T201    uint
	T202    uint
	T203    float64
	P501A   float64
	P501B   uint
	P502A   float64
	P502B   uint
	P503A   float64
	P503B   uint
	P504A   float64
	P504B   uint
	P505A   float64
	P505B   uint
	P506A   float64
	P506B   uint
	P507A   float64
	P507B   uint
	P508A   float64
	P508B   uint
	P510A   float64
	P510B   uint
	P511A   float64
	P511B   uint
	P512A   float64
	P512B   uint
	P521A   float64
	P521B   uint
	P522A   float64
	P522B   uint
	P523A   float64
	P523B   uint
	P524A   float64
	P524B   uint
	P538A   float64
	P538B   uint
	P539A   float64
	P539B   uint
	P546A   float64
	P546B   uint
	P548A   float64
	P548B   uint
	P549A   float64
	P549B   uint
	P550A   float64
	P550B   uint
	P551A   float64
	P551B   uint
	P572A   float64
	P572B   uint
	Pnl509  uint
	Pnn509  uint
	Pnl513  uint
	Pnn513  uint
	Pnl514  uint
	Pnn514  uint
	Pnl515  uint
	Pnn515  uint
	Pnl516  uint
	Pnn516  uint
	Pnl517  uint
	Pnn517  uint
	Pnl518  uint
	Pnn518  uint
	Pnl519  uint
	Pnn519  uint
	Pnl520  uint
	Pnn520  uint
	Pnl525  uint
	Pnn525  uint
	Pnl526  uint
	Pnn526  uint
	Pnl527  uint
	Pnn527  uint
	Pnl528  uint
	Pnn528  uint
	Pnl529  uint
	Pnn529  uint
	Pnl530  uint
	Pnn530  uint
	Pnl531  uint
	Pnn531  uint
	Pnl532  uint
	Pnn532  uint
	Pnl533  uint
	Pnn533  uint
	Pnl534  uint
	Pnn534  uint
	Pnl535  uint
	Pnn535  uint
	Pnl536  uint
	Pnn536  uint
	Pnl537  uint
	Pnn537  uint
	Pnl540  uint
	Pnn540  uint
	Pnl541  uint
	Pnn541  uint
	Pnl542  uint
	Pnn542  uint
	Pnl543  uint
	Pnn543  uint
	Pnl544  uint
	Pnn544  uint
	Pnl545  uint
	Pnn545  uint
	Pnl547  uint
	Pnn547  uint
	Pnl552  uint
	Pnn552  uint
	Pnl553  uint
	Pnn553  uint
	Pnl554  uint
	Pnn554  uint
	Pnl555  uint
	Pnn555  uint
	Pnl556  uint
	Pnn556  uint
	Pnl557  uint
	Pnn557  uint
	Pnl559  uint
	Pnn559  uint
	Pnl560  uint
	Pnn560  uint
	Pnl561  uint
	Pnn561  uint
	Pnl562  uint
	Pnn562  uint
	Pnl563  uint
	Pnn563  uint
	Pnl564  uint
	Pnn564  uint
	Pnl565  uint
	Pnn565  uint
	Pnl566  uint
	Pnn566  uint
	Pnl567  uint
	Pnn567  uint
	Pnl568  uint
	Pnn568  uint
	Pnl569  uint
	Pnn569  uint
	Pnl570  uint
	Pnn570  uint
	Pnl571  uint
	Pnn571  uint
	Pnl573  uint
	Pnn573  uint
	Pnl574  uint
	Pnn574  uint
	Pnl575  uint
	Pnn575  uint
	Pnl577  uint
	Pnn577  uint
	Pnl579  uint
	Pnn579  uint
	Pnl580  uint
	Pnn580  uint
	Pnl581  uint
	Pnn581  uint
	Pnl582  uint
	Pnn582  uint
	Pnl584  uint
	Pnn584  uint
	Pnl585  uint
	Pnn585  uint
}
