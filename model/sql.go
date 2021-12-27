package model

var index = map[string]string{
	"tally_sheet_no":                      `CREATE INDEX IF NOT EXISTS tally_sheet_no_index ON questionnaires (tally_sheet_no);`,
	"geo_code_id":                         `CREATE INDEX IF NOT EXISTS geocode_gist_idx ON ?TableName USING GIST (geo_code_id)`,
	"division":                            `CREATE INDEX IF NOT EXISTS division_index ON geo_codes (division);`,
	"district":                            `CREATE INDEX IF NOT EXISTS district_index ON geo_codes (district);`,
	"upazilla":                            `CREATE INDEX IF NOT EXISTS upazilla_index ON geo_codes (upazilla);`,
	"union":                               `CREATE INDEX IF NOT EXISTS union_index ON geo_codes ("union");`,
	"mouza":                               `CREATE INDEX IF NOT EXISTS mouza_index ON geo_codes (mouza);`,
	"mouza_search_idx":                    `CREATE INDEX IF NOT EXISTS mouza_search_idx ON geo_codes ((mouza || ' - ' || name_mouza) varchar_pattern_ops);`,
	"union_search_idx":                    `CREATE INDEX IF NOT EXISTS union_search_idx ON geo_codes (("union" || ' - ' || name_union) varchar_pattern_ops);`,
	"upazilla_search_idx":                 `CREATE INDEX IF NOT EXISTS upazilla_search_idx ON geo_codes ((upazilla || ' - ' || name_upazilla) varchar_pattern_ops);`,
	"district_search_isx":                 `CREATE INDEX IF NOT EXISTS district_search_idx ON geo_codes ((district || ' - ' || name_district) varchar_pattern_ops);`,
	"division_search_isx":                 `CREATE INDEX IF NOT EXISTS division_search_idx ON geo_codes ((division || ' - ' || name_division) varchar_pattern_ops);`,
	"questionnaires_booklet_number_index": `create index questionnaires_booklet_number_index on questionnaires (booklet_number);`,
	"questionnaires_geocode_id_index": `create index questionnaires_geocode_id_index
	on questionnaires (geocode_id);`,
	"overwritten_questionnaires_geocode_id_index": `create index overwritten_questionnaires_geocode_id_index
	on overwritten_questionnaires (geocode_id);`,
	"tally_sheets_geocode_id_index": `create index tally_sheets_geocode_id_index
	on tally_sheets (geocode_id);`,
}

var procedures = map[string]string{
	"update_tally_sheet_trigger": `
create or replace function update_tallysheet() returns trigger as
$update_tallysheet$
begin
    update tally_sheets
    set updated_buffalo_count         = (select sum((greatest(buffalo_at_farm, 0) + greatest(buffalo_at_home, 0)))
                                         from questionnaires
                                         where tally_sheet_no = OLD.tally_sheet_no),
        updated_cock_count            = (select sum((greatest(cock_at_farm, 0) + greatest(cock_at_home, 0)))
                                         from questionnaires
                                         where tally_sheet_no = OLD.tally_sheet_no),
        updated_cow_count             = (select sum((greatest(cow_at_farm, 0) + greatest(cow_at_home, 0)))
                                         from questionnaires
                                         where tally_sheet_no = OLD.tally_sheet_no),
        updated_duck_count            = (select sum((greatest(duck_at_farm, 0) + greatest(duck_at_home, 0)))
                                         from questionnaires
                                         where tally_sheet_no = OLD.tally_sheet_no),
        updated_goat_count            = (select sum((greatest(goat_at_farm, 0) + greatest(goat_at_home, 0)))
                                         from questionnaires
                                         where tally_sheet_no = OLD.tally_sheet_no),
        updated_sheep_count           = (select sum((greatest(sheep_at_farm, 0) + greatest(sheep_at_home, 0)))
                                         from questionnaires
                                         where tally_sheet_no = OLD.tally_sheet_no),
        updated_turkey_count          = (select sum((greatest(turkey_at_farm, 0) + greatest(turkey_at_home, 0)))
                                         from questionnaires
                                         where tally_sheet_no = OLD.tally_sheet_no),
        updated_agri_professionals    = (select count(q.questionnaire_num) as updated_agri_professionals
                                         from questionnaires q
                                         where q.tally_sheet_no = OLD.tally_sheet_no
                                           and q.agri_labor_code = 1),
        updated_fishing_professionals = (select count(q.questionnaire_num)
                                         from questionnaires q
                                         where q.tally_sheet_no = OLD.tally_sheet_no
                                           and q.is_fishing_related = 1),
        updated_house5_more           = (select count(q.questionnaire_num) as updated_house5_more
                                         from questionnaires q
                                         where q.tally_sheet_no = OLD.tally_sheet_no
                                           and q.cultivated_land > 0.05),
        updated_house_fisheries       = (select count(q.questionnaire_num) as updated_house_fisheries
                                         from questionnaires q
                                         where q.tally_sheet_no = OLD.tally_sheet_no
                                           and q.pond_land > 0),
        updated_house_no_land         = (select count(q.questionnaire_num) as updated_house_no_land
                                         from questionnaires q
                                         where q.tally_sheet_no = OLD.tally_sheet_no
                                           and (q.total_land = 0 or q.total_land is null)),
        updated_house_received_land   = (select count(q.questionnaire_num) as updated_house_received_land
                                         from questionnaires q
                                         where q.tally_sheet_no = OLD.tally_sheet_no
                                           and q.land_taken > 0),
        updated_total_house           = (select count(q.house_serial) as updated_total_house
                                         from questionnaires q
                                         where q.tally_sheet_no = OLD.tally_sheet_no
                                           and q.house_serial is not null
										   and (questionnaire_empty = false or questionnaire_empty is null))
    where tally_sheet_no = OLD.tally_sheet_no;
    return new;
end;
$update_tallysheet$ language plpgsql;
drop trigger if exists update_tallysheet on questionnaires;
create trigger update_tallysheet
    after update
    on questionnaires
    for each row
execute procedure update_tallysheet();`,
	"state_view": `
create or replace view stats as
with registered_today as
         (select count(b.number) as registered_today
          from booklets b
          where b.status = 'registered'
            and date(b.registered_on) >= current_date),
     registered_yesterday as
         (select count(b.number) as registered_yesterday
          from booklets b
          where b.status = 'registered'
            and date(b.registered_on) >= current_date - 2
            and date(b.registered_on) <= current_date - 1),
     inBatch_today as
         (select count(b.number) as inBatch_today
          from booklets b
          where b.status = 'inBatch'
            and date(b.added_in_batch_on) >= current_date),
     inBatch_yesterday as
         (select count(b.number) as inBatch_yesterday
          from booklets b
          where b.status = 'inBatch'
            and date(b.added_in_batch_on) >= current_date - 2
            and date(b.added_in_batch_on) <= current_date - 1),
     inCuttingStation_today as
         (select count(b.number) as inCuttingStation_today
          from booklets b
          where b.status = 'inCuttingStation'
            and date(b.cut_on) >= current_date),
     inCuttingStation_yesterday as
         (select count(b.number) as inCuttingStation_yesterday
          from booklets b
          where b.status = 'inCuttingStation'
            and date(b.cut_on) >= current_date - 2
            and date(b.cut_on) <= current_date - 1),
     inPreScanning_today as
         (select count(b.number) as inPreScanning_today
          from booklets b
          where b.status = 'inPreScanning'
            and date(b.prepared_on) >= current_date),
     inPreScanning_yesterday as
         (select count(b.number) as inPreScanning_yesterday
          from booklets b
          where b.status = 'inPreScanning'
            and date(b.prepared_on) >= current_date - 2
            and date(b.prepared_on) <= current_date - 1),
     inScanningStation_today as
         (select count(b.number) as inScanningStation_today
          from booklets b
          where b.status = 'inScanningStation'
            and date(b.scanned_on) >= current_date),
     inScanningStation_yesterday as
         (select count(b.number) as inScanningStation_yesterday
          from booklets b
          where b.status = 'inScanningStation'
            and date(b.scanned_on) >= current_date - 2
            and date(b.scanned_on) <= current_date - 1),
     archived_today as
         (select count(b.number) as archived_today
          from booklets b
          where b.status = 'archived'
            and date(b.archived_on) >= current_date),
     archived_yesterday as
         (select count(b.number) as archived_yesterday
          from booklets b
          where b.status = 'archived'
            and date(b.archived_on) >= current_date - 2
            and date(b.archived_on) <= current_date - 1)

select registered_today.registered_today,
       registered_yesterday.registered_yesterday,
       inBatch_today.inBatch_today,
       inBatch_yesterday.inBatch_yesterday,
       inCuttingStation_today.inCuttingStation_today,
       inCuttingStation_yesterday.inCuttingStation_yesterday,
       inPreScanning_today.inPreScanning_today,
       inPreScanning_yesterday.inPreScanning_yesterday,
       inScanningStation_today.inScanningStation_today,
       inScanningStation_yesterday.inScanningStation_yesterday,
       archived_today.archived_today,
       archived_yesterday.archived_yesterday
from registered_today,
     registered_yesterday,
     inBatch_today,
     inBatch_yesterday,
     inCuttingStation_today,
     inCuttingStation_yesterday,
     inPreScanning_today,
     inPreScanning_yesterday,
     inScanningStation_today,
     inScanningStation_yesterday,
     archived_today,
     archived_yesterday;`,
	"cascade_questionnaires_tally_sheet_no_fkey": `
	alter table questionnaires drop constraint questionnaires_tally_sheet_no_fkey;

alter table questionnaires
	add constraint questionnaires_tally_sheet_no_fkey
		foreign key (tally_sheet_no) references tally_sheets
			on update cascade;`,
	"cascade_questionnaires_booklet_number_fkey": `
	alter table questionnaires drop constraint questionnaires_booklet_number_fkey;

alter table questionnaires
	add constraint questionnaires_booklet_number_fkey
		foreign key (booklet_number) references booklets
			on update cascade;`,
	"cascade_overwritten_questionnaires_tally_sheet_no_fkey": `
	alter table overwritten_questionnaires drop constraint overwritten_questionnaires_tally_sheet_no_fkey;

alter table overwritten_questionnaires
	add constraint overwritten_questionnaires_tally_sheet_no_fkey
		foreign key (tally_sheet_no) references tally_sheets
			on update cascade;`,
	"cascade_overwritten_questionnaires_booklet_number_fkey": `
	alter table overwritten_questionnaires drop constraint overwritten_questionnaires_booklet_number_fkey;

alter table overwritten_questionnaires
	add constraint overwritten_questionnaires_booklet_number_fkey
		foreign key (booklet_number) references booklets
			on update cascade;
`,
	"cascade_events_archive_box_number_fkey": `
	alter table events drop constraint events_archive_box_number_fkey;

alter table events
	add constraint events_archive_box_number_fkey
		foreign key (archive_box_number) references archive_boxes
			on update cascade;
`,
	"cascade_events_booklet_number_fkey": `
	alter table events drop constraint events_booklet_number_fkey;

alter table events
	add constraint events_booklet_number_fkey
		foreign key (booklet_number) references booklets
			on update cascade;
`,
	"cascade_booklets_archive_box_number_fkey": `
	alter table booklets drop constraint booklets_archive_box_number_fkey;

alter table booklets
	add constraint booklets_archive_box_number_fkey
		foreign key (archive_box_number) references archive_boxes
			on update cascade;
`,
	"cascade_booklet_geocode_fkey": `alter table booklets drop constraint booklets_geo_code_id_fkey;

alter table booklets
	add constraint booklets_geo_code_id_fkey
		foreign key (geo_code_id) references geo_codes
			on update cascade;
`,
	"cascade_questionnaires_geocode_fkey": `alter table questionnaires drop constraint questionnaires_geocode_id_fkey;

alter table questionnaires
	add constraint questionnaires_geocode_id_fkey
		foreign key (geocode_id) references geo_codes
			on update cascade;`,
	"cascade_overwritten_questionnaire_geocode_fkey": `alter table overwritten_questionnaires drop constraint overwritten_questionnaires_geocode_id_fkey;

alter table overwritten_questionnaires
	add constraint overwritten_questionnaires_geocode_id_fkey
		foreign key (geocode_id) references geo_codes
			on update cascade;`,
	"cascade_tally_sheets_geocode_fkey": `alter table tally_sheets drop constraint tally_sheets_geocode_id_fkey;

alter table tally_sheets
	add constraint tally_sheets_geocode_id_fkey
		foreign key (geocode_id) references geo_codes
			on update cascade;
`,
}

type referenceData struct {
	Id   int
	Name string
}

var referencesData = map[string][]referenceData{
	"crop": {
		referenceData{Id: 101, Name: "Aush"},
		referenceData{Id: 102, Name: "Amond"},
		referenceData{Id: 103, Name: "Boro"},
		referenceData{Id: 104, Name: "Wheat"},
		referenceData{Id: 105, Name: "Maize"},
		referenceData{Id: 106, Name: "Kaon"},
		referenceData{Id: 107, Name: "Barley"},
		referenceData{Id: 108, Name: "Peanut"},
		referenceData{Id: 109, Name: "Pearl millet"},
		referenceData{Id: 110, Name: "Sorghum"},
		referenceData{Id: 111, Name: "Others Grains"},
		referenceData{Id: 112, Name: "Lentil"},
		referenceData{Id: 113, Name: "Saffran Pulse"},
		referenceData{Id: 114, Name: "Moog Pulse"},
		referenceData{Id: 115, Name: "Vigna mungo"},
		referenceData{Id: 116, Name: "Pea"},
		referenceData{Id: 117, Name: "Gram"},
		referenceData{Id: 118, Name: "Pigeon Pea"},
		referenceData{Id: 119, Name: "Fallon"},
		referenceData{Id: 120, Name: "Others Pulse"},
		referenceData{Id: 121, Name: "Potato"},
		referenceData{Id: 122, Name: "Brinjal"},
		referenceData{Id: 123, Name: "Radish"},
		referenceData{Id: 124, Name: "Bean"},
		referenceData{Id: 125, Name: "Tomato"},
		referenceData{Id: 126, Name: "Snake Gourd"},
		referenceData{Id: 127, Name: "Multitude"},
		referenceData{Id: 128, Name: "Ladies Finger"},
		referenceData{Id: 129, Name: "Cucumber"},
		referenceData{Id: 130, Name: "Bitter Gourd / Momordica / Charantia"},
		referenceData{Id: 131, Name: "Gourd"},
		referenceData{Id: 132, Name: "Pumplin"},
		referenceData{Id: 133, Name: "Wax gourd"},
		referenceData{Id: 134, Name: "Cauliflower"},
		referenceData{Id: 135, Name: "Cabbage"},
		referenceData{Id: 136, Name: "Broccoli"},
		referenceData{Id: 137, Name: "Momordica dioica"},
		referenceData{Id: 138, Name: "Sweet Potato"},
		referenceData{Id: 139, Name: "Daata"},
		referenceData{Id: 140, Name: "Taro"},
		referenceData{Id: 141, Name: "Long Beans"},
		referenceData{Id: 142, Name: "Luffa acutangula"},
		referenceData{Id: 143, Name: "Corrot"},
		referenceData{Id: 144, Name: "Kholrabi"},
		referenceData{Id: 145, Name: "Turnip"},
		referenceData{Id: 146, Name: "Cucumber"},
		referenceData{Id: 147, Name: "Capsicum"},
		referenceData{Id: 148, Name: "Sponge gourd"},
		referenceData{Id: 149, Name: "Beetroot"},
		referenceData{Id: 150, Name: "Others veg"},
		referenceData{Id: 151, Name: "Red Spinach"},
		referenceData{Id: 152, Name: "Malabar spinach"},
		referenceData{Id: 153, Name: "Palong shak"},
		referenceData{Id: 154, Name: "Mint Leaves"},
		referenceData{Id: 155, Name: "Lettuce"},
		referenceData{Id: 156, Name: "Other shak"},
		referenceData{Id: 157, Name: "Onion"},
		referenceData{Id: 158, Name: "Garlic"},
		referenceData{Id: 159, Name: "Ginger"},
		referenceData{Id: 160, Name: "Turmeric"},
		referenceData{Id: 161, Name: "Chili"},
		referenceData{Id: 162, Name: "Coriander"},
		referenceData{Id: 163, Name: "Black Cumin"},
		referenceData{Id: 164, Name: "Fennil seed"},
		referenceData{Id: 165, Name: "Cumin"},
		referenceData{Id: 166, Name: "Others Spices"},
		referenceData{Id: 167, Name: "Mustard"},
		referenceData{Id: 168, Name: "Soybeans"},
		referenceData{Id: 169, Name: "Nuts"},
		referenceData{Id: 170, Name: "Sesame"},
		referenceData{Id: 171, Name: "Linseed"},
		referenceData{Id: 172, Name: "Sunflower"},
		referenceData{Id: 173, Name: "Verenda"},
		referenceData{Id: 174, Name: "others oil seeds"},
		referenceData{Id: 175, Name: "Banana"},
		referenceData{Id: 176, Name: "Papaya"},
		referenceData{Id: 177, Name: "Water Melon"},
		referenceData{Id: 178, Name: "Pronunciation"},
		referenceData{Id: 179, Name: "Pine Apple"},
		referenceData{Id: 180, Name: "Strawberry"},
		referenceData{Id: 181, Name: "Others fruit"},
		referenceData{Id: 182, Name: "Jute"},
		referenceData{Id: 183, Name: "Corpus cotton"},
		referenceData{Id: 184, Name: "Others"},
		referenceData{Id: 185, Name: "Sugar Cane"},
		referenceData{Id: 186, Name: "Others"},
		referenceData{Id: 187, Name: "Tobacco"},
		referenceData{Id: 188, Name: "Others"},
		referenceData{Id: 189, Name: "Aloevera"},
		referenceData{Id: 190, Name: "Others"},
		referenceData{Id: 191, Name: "Nightshade"},
		referenceData{Id: 192, Name: "Marigold"},
		referenceData{Id: 193, Name: "Chandramallika"},
		referenceData{Id: 194, Name: "Dahlia"},
		referenceData{Id: 195, Name: "Gladioilus"},
		referenceData{Id: 196, Name: "Gerbera"},
		referenceData{Id: 197, Name: "Others"},
		referenceData{Id: 198, Name: "Chan"},
		referenceData{Id: 199, Name: "Sesbania bispinosa"},
		referenceData{Id: 200, Name: "Other fuel"},
		referenceData{Id: 201, Name: "Napier grass"},
		referenceData{Id: 202, Name: "other cow food"},
		referenceData{Id: 203, Name: "Seeded"},
	},
	"tree": {
		referenceData{Id: 501, Name: "Mango"},
		referenceData{Id: 502, Name: "Berry"},
		referenceData{Id: 503, Name: "Jack Fruit"},
		referenceData{Id: 504, Name: "Litchi"},
		referenceData{Id: 505, Name: "Guava"},
		referenceData{Id: 506, Name: "Coconut"},
		referenceData{Id: 507, Name: "Plum"},
		referenceData{Id: 508, Name: "Hog Plum"},
		referenceData{Id: 509, Name: "Olives"},
		referenceData{Id: 510, Name: "Date"},
		referenceData{Id: 511, Name: "Palmyra"},
		referenceData{Id: 512, Name: "Bell"},
		referenceData{Id: 513, Name: "wood Apple"},
		referenceData{Id: 514, Name: "Jamrul"},
		referenceData{Id: 515, Name: "Koromcha"},
		referenceData{Id: 516, Name: "Graviola"},
		referenceData{Id: 517, Name: "Custurd apple"},
		referenceData{Id: 518, Name: "Pomcgranate"},
		referenceData{Id: 519, Name: "Sapodilla"},
		referenceData{Id: 520, Name: "Monkey jack"},
		referenceData{Id: 521, Name: "Averrhoa Carambola"},
		referenceData{Id: 522, Name: "Tamarind"},
		referenceData{Id: 523, Name: "Lemon"},
		referenceData{Id: 524, Name: "Grapefruit"},
		referenceData{Id: 525, Name: "Amla"},
		referenceData{Id: 526, Name: "Lotkon"},
		referenceData{Id: 527, Name: "Phyllanthin acid"},
		referenceData{Id: 528, Name: "Karambel"},
		referenceData{Id: 529, Name: "Orange"},
		referenceData{Id: 530, Name: "Citrus macioptera"},
		referenceData{Id: 531, Name: "Malta"},
		referenceData{Id: 532, Name: "Averrhoa billimbi"},
		referenceData{Id: 533, Name: "Gab"},
		referenceData{Id: 534, Name: "Fig"},
		referenceData{Id: 535, Name: "Dragon"},
		referenceData{Id: 536, Name: "Rambutan"},
		referenceData{Id: 537, Name: "Others"},
		referenceData{Id: 538, Name: "Battle Leaf"},
		referenceData{Id: 539, Name: "Areca Catechu"},
		referenceData{Id: 540, Name: "Tea"},
		referenceData{Id: 541, Name: "Others"},
		referenceData{Id: 542, Name: "Bay leaves"},
		referenceData{Id: 543, Name: "Cinnamon"},
		referenceData{Id: 544, Name: "Cardamom"},
		referenceData{Id: 545, Name: "Others"},
		referenceData{Id: 546, Name: "Bamboo"},
		referenceData{Id: 547, Name: "Cane"},
		referenceData{Id: 548, Name: "Hardwood Tree"},
		referenceData{Id: 549, Name: "Rain Tree"},
		referenceData{Id: 550, Name: "Mahogany"},
		referenceData{Id: 551, Name: "Tectona Grandis"},
		referenceData{Id: 552, Name: "Akashmoni"},
		referenceData{Id: 553, Name: "Eucalyptus"},
		referenceData{Id: 554, Name: "Sisu"},
		referenceData{Id: 555, Name: "Gamari"},
		referenceData{Id: 556, Name: "Garjon"},
		referenceData{Id: 557, Name: "Banayan"},
		referenceData{Id: 558, Name: "Schumannianthus dichotomus"},
		referenceData{Id: 559, Name: "Southern cattail"},
		referenceData{Id: 560, Name: "Monoon longifolium"},
		referenceData{Id: 561, Name: "Devil's Tree"},
		referenceData{Id: 562, Name: "Jarul tree"},
		referenceData{Id: 563, Name: "Babla"},
		referenceData{Id: 564, Name: "Chambula"},
		referenceData{Id: 565, Name: "Golden Shower"},
		referenceData{Id: 566, Name: "Ficus elastica"},
		referenceData{Id: 567, Name: "Agar wood"},
		referenceData{Id: 568, Name: "Shorea robusta"},
		referenceData{Id: 569, Name: "Artocarpus"},
		referenceData{Id: 570, Name: "Albizia lebback"},
		referenceData{Id: 571, Name: "Others"},
		referenceData{Id: 572, Name: "Neem"},
		referenceData{Id: 573, Name: "Terminalia arjuna"},
		referenceData{Id: 574, Name: "Terminalia bellirica"},
		referenceData{Id: 575, Name: "Terminalia Chebula"},
		referenceData{Id: 576, Name: "Others"},
		referenceData{Id: 577, Name: "Shojina"},
		referenceData{Id: 578, Name: "Others"},
		referenceData{Id: 579, Name: "Rose"},
		referenceData{Id: 580, Name: "Mimusops Elengi"},
		referenceData{Id: 581, Name: "Neolamarckia cadamba"},
		referenceData{Id: 582, Name: "Delonix regia"},
		referenceData{Id: 583, Name: "Others"},
		referenceData{Id: 584, Name: "Shimul cotton"},
		referenceData{Id: 585, Name: "Morus"},
		referenceData{Id: 586, Name: "Others"},
	},
}
