package routes

type Field struct {
	label     string
	name      string
	value     string
	pattern   string
	inputMask string
}

var table9 = [][]Field{
	{
		{"01 Total member of house", "TotalMemberMale", "TotalMemberMale", "\\d{5}", "'mask' : '99999'"},
		{"01 Total member of house", "TotalMemberFemale", "TotalMemberFemale", "\\d{5}", "'mask' : '99999'"},
		{"01 Total member of house", "TotalMemberHijra", "TotalMemberHijra", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"02 People between 10 to 14yo", "PeopleBetween10to14Male", "PeopleBetween10to14Male", "\\d{5}", "'mask' : '99999'"},
		{"02 People between 10 to 14yo", "PeopleBetween10to14Female", "PeopleBetween10to14Female", "\\d{5}", "'mask' : '99999'"},
		{"02 People between 10 to 14yo", "PeopleBetween10to14Hijra", "PeopleBetween10to14Hijra", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"03 People above 15", "PeopleAbove15Male", "PeopleAbove15Male", "\\d{5}", "'mask' : '99999'"},
		{"03 People above 15", "PeopleAbove15Female", "PeopleAbove15Female", "\\d{5}", "'mask' : '99999'"},
		{"03 People above 15", "PeopleAbove15Hijra", "PeopleAbove15Hijra", "\\d{5}", "'mask' : '99999'"},
	},
}

var table10 = []Field{
	{"04 Total own land", "TotalLand", "TotalLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"05 Land given", "LandGiven", "LandGiven", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"06 Land taken", "LandTaken", "LandTaken", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"07 Total operating Land", "OperatingLand", "OperatingLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"08 Land  for residence", "ResidenceLand", "ResidenceLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"09 Land used for business entity", "BusinessEntityLand", "BusinessEntityLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"10 Land in sink-channel, bush etc", "SinkChannelLand", "SinkChannelLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"11 Permanent unused Land", "PermanentUnusedLand", "PermanentUnusedLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"12 Uncultivated Land", "UncultivatedLand", "UncultivatedLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"13 Temporary cultivated Land", "TempCultivatedLand", "TempCultivatedLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"14 Permanent cultivated land (forest, tree and fruit)", "PermanentCultivatedLand", "PermanentCultivatedLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"15 Land under pond/lake (without bank)", "PondLandNotBlank", "PondLandNotBlank", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"16 Land under nursery", "NurseryLand", "NurseryLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"17 Recent unused Land", "RecentUnusedLand", "RecentUnusedLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"18 Cultivated Land", "CultivatedLand", "CultivatedLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"19 Land under irrigation", "IrrigationLand", "IrrigationLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"20 Land under salt cultivation", "SaltCultivationLand", "SaltCultivationLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
}

var table11 = []Field{
	{"21 Land under Pond/Lake", "PondLand", "PondLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"22 Land under Fish cultivation", "FishCultivationLand", "FishCultivationLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"23 Land under paddy and fish cultivation", "PaddyCultivationLand", "PaddyCultivationLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"24 Land under mixed cultivation duck-hen and fish", "MixedCultivationLand", "MixedCultivationLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"25 Fish cultivation on salt cultivating land", "FishSaltCultiveLand", "FishSaltCultiveLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"26 Land under fish cultivation in cage", "FishCageCultiveLand", "FishCageCultiveLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"27 Land under creek", "CreekLand", "CreekLand", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
}

var table12 = [][]Field{
	{
		{"28 Cock - hen", "CockAtHome", "CockAtHome", "\\d{5}", "'mask' : '99999'"},
		{"28 Cock - hen", "CockAtFarm", "CockAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"29 Duck", "DuckAtHome", "DuckAtHome", "\\d{5}", "'mask' : '99999'"},
		{"29 Duck", "DuckAtFarm", "DuckAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"30 Pigeon", "PigeonAtHome", "PigeonAtHome", "\\d{5}", "'mask' : '99999'"},
		{"30 Pigeon", "PigeonAtFarm", "PigeonAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"31 Quail", "QuailAtHome", "QuailAtHome", "\\d{5}", "'mask' : '99999'"},
		{"31 Quail", "QuailAtFarm", "QuailAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"32 Turkey", "TurkeyAtHome", "TurkeyAtHome", "\\d{5}", "'mask' : '99999'"},
		{"32 Turkey", "TurkeyAtFarm", "TurkeyAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
}

var table13 = [][]Field{
	{
		{"33 Cow", "CowAtHome", "CowAtHome", "\\d{5}", "'mask' : '99999'"},
		{"33 Cow", "CowAtFarm", "CowAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"34 Buffalo", "BuffaloAtHome", "BuffaloAtHome", "\\d{5}", "'mask' : '99999'"},
		{"34 Buffalo", "BuffaloAtFarm", "BuffaloAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"35 Goat", "GoatAtHome", "GoatAtHome", "\\d{5}", "'mask' : '99999'"},
		{"35 Goat", "GoatAtFarm", "GoatAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"36 Sheep", "SheepAtHome", "SheepAtHome", "\\d{5}", "'mask' : '99999'"},
		{"36 Sheep", "SheepAtFarm", "SheepAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"37 Pig", "PigAtHome", "PigAtHome", "\\d{5}", "'mask' : '99999'"},
		{"37 Pig", "PigAtFarm", "PigAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"38 Horse", "HorseAtHome", "HorseAtHome", "\\d{5}", "'mask' : '99999'"},
		{"38 Horse", "HorseAtFarm", "HorseAtFarm", "\\d{5}", "'mask' : '99999'"},
	},
}

var table14Additional = [][]Field{
	{
		{"1 Additional cultivating", "TempAddCrop01Name", "TempAddCrop01Name", ".*", ""},
		{"1 Additional cultivating", "TempAddCrop01Id", "TempAddCrop01Id", "\\d{5}", "'mask' : '99999'"},
		{"1 Additional cultivating", "TempAddCrop01Surface", "TempAddCrop01Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"2 Additional cultivating", "TempAddCrop02Name", "TempAddCrop02Name", ".*", ""},
		{"2 Additional cultivating", "TempAddCrop02Id", "TempAddCrop02Id", "\\d{5}", "'mask' : '99999'"},
		{"2 Additional cultivating", "TempAddCrop02Surface", "TempAddCrop02Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"3 Additional cultivating", "TempAddCrop03Name", "TempAddCrop03Name", ".*", ""},
		{"3 Additional cultivating", "TempAddCrop03Id", "TempAddCrop03Id", "\\d{5}", "'mask' : '99999'"},
		{"3 Additional cultivating", "TempAddCrop03Surface", "TempAddCrop03Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"4 Additional cultivating", "TempAddCrop04Name", "TempAddCrop04Name", ".*", ""},
		{"4 Additional cultivating", "TempAddCrop04Id", "TempAddCrop04Id", "\\d{5}", "'mask' : '99999'"},
		{"4 Additional cultivating", "TempAddCrop04Surface", "TempAddCrop04Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"5 Additional cultivating", "TempAddCrop05Name", "TempAddCrop05Name", ".*", ""},
		{"5 Additional cultivating", "TempAddCrop05Id", "TempAddCrop05Id", "\\d{5}", "'mask' : '99999'"},
		{"5 Additional cultivating", "TempAddCrop05Surface", "TempAddCrop05Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"6 Additional cultivating", "TempAddCrop06Name", "TempAddCrop06Name", ".*", ""},
		{"6 Additional cultivating", "TempAddCrop06Id", "TempAddCrop06Id", "\\d{5}", "'mask' : '99999'"},
		{"6 Additional cultivating", "TempAddCrop06Surface", "TempAddCrop06Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"7 Additional cultivating", "TempAddCrop07Name", "TempAddCrop07Name", ".*", ""},
		{"7 Additional cultivating", "TempAddCrop07Id", "TempAddCrop07Id", "\\d{5}", "'mask' : '99999'"},
		{"7 Additional cultivating", "TempAddCrop07Surface", "TempAddCrop07Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"8 Additional cultivating", "TempAddCrop08Name", "TempAddCrop08Name", ".*", ""},
		{"8 Additional cultivating", "TempAddCrop08Id", "TempAddCrop08Id", "\\d{5}", "'mask' : '99999'"},
		{"8 Additional cultivating", "TempAddCrop08Surface", "TempAddCrop08Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"9 Additional cultivating", "TempAddCrop09Name", "TempAddCrop09Name", ".*", ""},
		{"9 Additional cultivating", "TempAddCrop09Id", "TempAddCrop09Id", "\\d{5}", "'mask' : '99999'"},
		{"9 Additional cultivating", "TempAddCrop09Surface", "TempAddCrop09Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
}

var table14R = []Field{
	{"101 Aush", "TempCrop101Aush", "TempCrop101Aush", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"102 Amond", "TempCrop102Almond", "TempCrop102Almond", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"103 Boro", "TempCrop103Boro", "TempCrop103Boro", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"104 Wheat", "TempCrop104Wheat", "TempCrop104Wheat", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"105 Maize", "TempCrop105Maize", "TempCrop105Maize", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"112 Lentil", "TempCrop112Lentil", "TempCrop112Lentil", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"113 Safran Pulse", "TempCrop113Safran", "TempCrop113Safran", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"114 Moog Pulse", "TempCrop114Moog", "TempCrop114Moog", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"121 Potato", "TempCrop121Potato", "TempCrop121Potato", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"122 Brinjal", "TempCrop122Brinjal", "TempCrop122Brinjal", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"123 Radish", "TempCrop123Radish", "TempCrop123Radish", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"124 Bean", "TempCrop124Bean", "TempCrop124Bean", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"125 Tomato", "TempCrop125Tomato", "TempCrop125Tomato", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"127 Multitude", "TempCrop127Multitude", "TempCrop127Multitude", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"128 Ladies Finger", "TempCrop128ladyfinger", "TempCrop128ladyfinger", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"129 Cucumber", "TempCrop129Cucumber", "TempCrop129Cucumber", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"130 Bitter Gourd / Momordica / Charantia", "TempCrop130BitterGourd", "TempCrop130BitterGourd", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
}

var table14L = []Field{
	{"131 Gourd", "TempCrop131Ground", "TempCrop131Ground", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"132 Pumpkin", "TempCrop132Pumplin", "TempCrop132Pumplin", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"134 Cauliflower", "TempCrop134Cauliflower", "TempCrop134Cauliflower", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"135 Cabbage", "TempCrop135Cabbage", "TempCrop135Cabbage", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"157 Onion", "TempCrop157Onion", "TempCrop157Onion", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"158 Garlic", "TempCrop158Garlic", "TempCrop158Garlic", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"159 Ginger", "TempCrop159Ginger", "TempCrop159Ginger", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"160 Turmeric", "TempCrop160Turmeric", "TempCrop160Turmeric", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"161 Chili", "TempCrop161Chili", "TempCrop161Chili", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"167 Mustard", "TempCrop167Mustard", "TempCrop167Mustard", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"169 Nuts", "TempCrop169Nuts", "TempCrop169Nuts", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"175 Banana", "TempCrop175Banana", "TempCrop175Banana", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"176 Papaya", "TempCrop176Papaya", "TempCrop176Papaya", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"177 Water Melon", "TempCrop177WMelon", "TempCrop177WMelon", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"179 Pine Apple", "TempCrop179PineApple", "TempCrop179PineApple", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"182 Jute", "TempCrop182Jute", "TempCrop182Jute", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"185 Sugar Cane", "TempCrop185SugarCane", "TempCrop185SugarCane", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	{"203 Seeded", "TempCrop203Seeded", "TempCrop203Seeded", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
}

var table15 = [][]Field{
	{
		{"501 Mango", "FixedCrop501MongoSurface", "FixedCrop501MongoSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"501 Mango", "FixedCrop501MongoCount", "FixedCrop501MongoCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"502 Berry", "FixedCrop502BerrySurface", "FixedCrop502BerrySurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"502 Berry", "FixedCrop502BerryCount", "FixedCrop502BerryCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"503 Jack Fruit", "FixedCrop503JackFruitSurface", "FixedCrop503JackFruitSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"503 Jack Fruit", "FixedCrop503JackFruitCount", "FixedCrop503JackFruitCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"504 Litchi", "FixedCrop504LitchiSurface", "FixedCrop504LitchiSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"504 Litchi", "FixedCrop504LitchiCount", "FixedCrop504LitchiCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"505 Guava", "FixedCrop505GuavaSurface", "FixedCrop505GuavaSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"505 Guava", "FixedCrop505GuavaCount", "FixedCrop505GuavaCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"506 Coconut", "FixedCrop506CoconutSurface", "FixedCrop506CoconutSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"506 Coconut", "FixedCrop506CoconutCount", "FixedCrop506CoconutCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"507 Plum", "FixedCrop507PlumSurface", "FixedCrop507PlumSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"507 Plum", "FixedCrop507PlumCount", "FixedCrop507PlumCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"508 Hog Plum", "FixedCrop508HogPlumSurface", "FixedCrop508HogPlumSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"508 Hog Plum", "FixedCrop508HogPlumCount", "FixedCrop508HogPlumCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"510 Date", "FixedCrop510DateSurface", "FixedCrop510DateSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"510 Date", "FixedCrop510DateCount", "FixedCrop510DateCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"511 Palmyra", "FixedCrop511PalmyraSurface", "FixedCrop511PalmyraSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"511 Palmyra", "FixedCrop511PalmyraCount", "FixedCrop511PalmyraCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"512 Bell", "FixedCrop512BellSurface", "FixedCrop512BellSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"512 Bell", "FixedCrop512BellCount", "FixedCrop512BellCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"521 Averrhoa Carambola", "FixedCrop521AverrhoaSurface", "FixedCrop521AverrhoaSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"521 Averrhoa Carambola", "FixedCrop521AverrhoaCount", "FixedCrop521AverrhoaCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"522 Tamarind", "FixedCrop522TamarindSurface", "FixedCrop522TamarindSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"522 Tamarind", "FixedCrop522TamarindCount", "FixedCrop522TamarindCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"523 Lemon", "FixedCrop523LemonSurface", "FixedCrop523LemonSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"523 Lemon", "FixedCrop523LemonCount", "FixedCrop523LemonCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"524 Grapefruit", "FixedCrop524GrapeSurface", "FixedCrop524GrapeSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"524 Grapefruit", "FixedCrop524GrapeCount", "FixedCrop524GrapeCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"538 Battle Leaf", "FixedCrop538BattleLeafSurface", "FixedCrop538BattleLeafSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"538 Battle Leaf", "FixedCrop538BattleLeafCount", "FixedCrop538BattleLeafCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"539 Areca Catechu", "FixedCrop539ArecaCatechuSurface", "FixedCrop539ArecaCatechuSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"539 Areca Catechu", "FixedCrop539ArecaCatechuCount", "FixedCrop539ArecaCatechuCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"546 Bamboo", "FixedCrop546BambooSurface", "FixedCrop546BambooSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"546 Bamboo", "FixedCrop546BambooCount", "FixedCrop546BambooCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"548 Hardwood Tree", "FixedCrop548HardwoodSurface", "FixedCrop548HardwoodSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"548 Hardwood Tree", "FixedCrop548HardwoodCount", "FixedCrop548HardwoodCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"549 Rain Tree", "FixedCrop549RainTreeSurface", "FixedCrop549RainTreeSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"549 Rain Tree", "FixedCrop549RainTreeCount", "FixedCrop549RainTreeCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"550 Mahogany", "FixedCrop550MahoganySurface", "FixedCrop550MahoganySurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"550 Mahogany", "FixedCrop550MahoganyCount", "FixedCrop550MahoganyCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"551 Tectona Grandis", "FixedCrop551TectonaGrandisSurface", "FixedCrop551TectonaGrandisSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"551 Tectona Grandis", "FixedCrop551TectonaGrandisCount", "FixedCrop551TectonaGrandisCount", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"572 Neem", "FixedCrop572NeemSurface", "FixedCrop572NeemSurface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
		{"572 Neem", "FixedCrop572NeemCount", "FixedCrop572NeemCount", "\\d{5}", "'mask' : '99999'"},
	},
}

var table15Additional = [][]Field{
	{
		{"1 Additional fixed cultivating", "FixedCropAddTree01Name", "FixedCropAddTree01Name", ".*", ""},
		{"1 Additional fixed cultivating", "FixedCropAddTree01Id", "FixedCropAddTree01Id", "\\d{5}", "'mask' : '99999'"},
		{"1 Additional fixed cultivating", "FixedCropAddTree01Count", "FixedCropAddTree01Count", "\\d{5}", "'mask' : '99999'"},
		{"1 Additional fixed cultivating", "FixedCropAddTree01Surface", "FixedCropAddTree01Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"2 Additional fixed cultivating", "FixedCropAddTree02Name", "FixedCropAddTree02Name", ".*", ""},
		{"2 Additional fixed cultivating", "FixedCropAddTree02Id", "FixedCropAddTree02Id", "\\d{5}", "'mask' : '99999'"},
		{"2 Additional fixed cultivating", "FixedCropAddTree02Count", "FixedCropAddTree02Count", "\\d{5}", "'mask' : '99999'"},
		{"2 Additional fixed cultivating", "FixedCropAddTree02Surface", "FixedCropAddTree02Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"3 Additional fixed cultivating", "FixedCropAddTree03Name", "FixedCropAddTree03Name", ".*", ""},
		{"3 Additional fixed cultivating", "FixedCropAddTree03Id", "FixedCropAddTree03Id", "\\d{5}", "'mask' : '99999'"},
		{"3 Additional fixed cultivating", "FixedCropAddTree03Count", "FixedCropAddTree03Count", "\\d{5}", "'mask' : '99999'"},
		{"3 Additional fixed cultivating", "FixedCropAddTree03Surface", "FixedCropAddTree03Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"4 Additional fixed cultivating", "FixedCropAddTree04Name", "FixedCropAddTree04Name", ".*", ""},
		{"4 Additional fixed cultivating", "FixedCropAddTree04Id", "FixedCropAddTree04Id", "\\d{5}", "'mask' : '99999'"},
		{"4 Additional fixed cultivating", "FixedCropAddTree04Count", "FixedCropAddTree04Count", "\\d{5}", "'mask' : '99999'"},
		{"4 Additional fixed cultivating", "FixedCropAddTree04Surface", "FixedCropAddTree04Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"5 Additional fixed cultivating", "FixedCropAddTree05Name", "FixedCropAddTree05Name", ".*", ""},
		{"5 Additional fixed cultivating", "FixedCropAddTree05Id", "FixedCropAddTree05Id", "\\d{5}", "'mask' : '99999'"},
		{"5 Additional fixed cultivating", "FixedCropAddTree05Count", "FixedCropAddTree05Count", "\\d{5}", "'mask' : '99999'"},
		{"5 Additional fixed cultivating", "FixedCropAddTree05Surface", "FixedCropAddTree05Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"6 Additional fixed cultivating", "FixedCropAddTree06Name", "FixedCropAddTree06Name", ".*", ""},
		{"6 Additional fixed cultivating", "FixedCropAddTree06Id", "FixedCropAddTree06Id", "\\d{5}", "'mask' : '99999'"},
		{"6 Additional fixed cultivating", "FixedCropAddTree06Count", "FixedCropAddTree06Count", "\\d{5}", "'mask' : '99999'"},
		{"6 Additional fixed cultivating", "FixedCropAddTree06Surface", "FixedCropAddTree06Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"7 Additional fixed cultivating", "FixedCropAddTree07Name", "FixedCropAddTree07Name", ".*", ""},
		{"7 Additional fixed cultivating", "FixedCropAddTree07Id", "FixedCropAddTree07Id", "\\d{5}", "'mask' : '99999'"},
		{"7 Additional fixed cultivating", "FixedCropAddTree07Count", "FixedCropAddTree07Count", "\\d{5}", "'mask' : '99999'"},
		{"7 Additional fixed cultivating", "FixedCropAddTree07Surface", "FixedCropAddTree07Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
	{
		{"8 Additional fixed cultivating", "FixedCropAddTree08Name", "FixedCropAddTree08Name", ".*", ""},
		{"8 Additional fixed cultivating", "FixedCropAddTree08Id", "FixedCropAddTree08Id", "\\d{5}", "'mask' : '99999'"},
		{"8 Additional fixed cultivating", "FixedCropAddTree08Count", "FixedCropAddTree08Count", "\\d{5}", "'mask' : '99999'"},
		{"8 Additional fixed cultivating", "FixedCropAddTree08Surface", "FixedCropAddTree08Surface", "\\d{2}.\\d{2}", "'mask' : '99.99'"},
	},
}

var table16 = [][]Field{
	{
		{"39 Tractor", "TractorManual", "", "\\d{5}", "'mask' : '99999'"},
		{"39 Tractor", "TractorDiesel", "TractorDiesel", "\\d{5}", "'mask' : '99999'"},
		{"39 Tractor", "TractorElectric", "", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"40 Power Tiller", "PowerTrillerManual", "", "\\d{5}", "'mask' : '99999'"},
		{"40 Power Tiller", "PowerTrillerDiesel", "PowerTrillerDiesel", "\\d{5}", "'mask' : '99999'"},
		{"40 Power Tiller", "PowerTrillerElectric", "", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"41 Power Pump", "PowerPumpManual", "", "\\d{5}", "'mask' : '99999'"},
		{"41 Power Pump", "PowerPumpDiesel", "PowerPumpDiesel", "\\d{5}", "'mask' : '99999'"},
		{"41 Power Pump", "PowerPumpElectric", "PowerPumpElectric", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"42 Deep/Shallow Tube Well", "TubeWellManual", "", "\\d{5}", "'mask' : '99999'"},
		{"42 Deep/Shallow Tube Well", "TubeWellDiesel", "TubeWellDiesel", "\\d{5}", "'mask' : '99999'"},
		{"42 Deep/Shallow Tube Well", "TubeWellElectric", "TubeWellElectric", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"43 Planting / Sowing machine", "SowingMachineManual", "SowingMachineManual", "\\d{5}", "'mask' : '99999'"},
		{"43 Planting / Sowing machine", "SowingMachineDiesel", "SowingMachineDiesel", "\\d{5}", "'mask' : '99999'"},
		{"43 Planting / Sowing machine", "SowingMachineElectric", "", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"44 Crop Cutting Machine", "CuttingMachineManual", "CuttingMachineManual", "\\d{5}", "'mask' : '99999'"},
		{"44 Crop Cutting Machine", "CuttingMachineDiesel", "CuttingMachineDiesel", "\\d{5}", "'mask' : '99999'"},
		{"44 Crop Cutting Machine", "CuttingMachineElectric", "", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"45 Harvest Machine", "HarvestMachineManual", "HarvestMachineManual", "\\d{5}", "'mask' : '99999'"},
		{"45 Harvest Machine", "HarvestMachineDiesel", "HarvestMachineDiesel", "\\d{5}", "'mask' : '99999'"},
		{"45 Harvest Machine", "HarvestMachineElectric", "HarvestMachineElectric", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"46 Fertilizer Applying Machine", "FertilizerManual", "FertilizerManual", "\\d{5}", "'mask' : '99999'"},
		{"46 Fertilizer Applying Machine", "FertilizerDiesel", "FertilizerDiesel", "\\d{5}", "'mask' : '99999'"},
		{"46 Fertilizer Applying Machine", "FertilizerElectric", "", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"47 Fishing Boat / Trawler", "LinerManual", "LinerManual", "\\d{5}", "'mask' : '99999'"},
		{"47 Fishing Boat / Trawler", "LinerDiesel", "LinerDiesel", "\\d{5}", "'mask' : '99999'"},
		{"47 Fishing Boat / Trawler", "LinerElectric", "", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"48 Fishing Net (commercial)", "FishingNetManual", "FishingNetManual", "\\d{5}", "'mask' : '99999'"},
		{"48 Fishing Net (commercial)", "FishingNetDiesel", "", "\\d{5}", "'mask' : '99999'"},
		{"48 Fishing Net (commercial)", "FishingNetElectric", "", "\\d{5}", "'mask' : '99999'"},
	},
	{
		{"49 Plow", "PlowManual", "PlowManual", "\\d{5}", "'mask' : '99999'"},
		{"49 Plow", "PlowDiesel", "", "\\d{5}", "'mask' : '99999'"},
		{"49 Plow", "PlowElectric", "", "\\d{5}", "'mask' : '99999'"},
	},
}
