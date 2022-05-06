package routes

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (srv *Server) FormatHouseholdPermanentCrops(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	c, err := srv.Db.GetTemporaryCrops(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	p := message.NewPrinter(language.English)

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> </small></h4>
	<h5>Data for table number : %s</h5>
	<h5>Number of farm holdings : %s</h5>
	<h5>Crop area : %s</h5>
	<table id="datatable-buttons" class="table table-striped">
		<thead>
			<tr>
				<th>Crop code and name</th>
				<th>Total temporary crop area</th>
				<th>Percentage of crop area</th>
			</tr>
		</thead>
		<tbody>
			<tr> <td><b>501 Mango</b> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>502 Berry</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>503 Jack Fruit</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>504 WLitchiheat</b></td> <td>%.2f</td> <td>%s</td> </tr>         
			<tr> <td><b>505 Guava</b></td> <td>%.2f</td> <td>%s</td> </tr>         
			<tr> <td><b>506 Coconut millet</b></td> <td>%.2f</td> <td>%s</td> </tr>
			<tr> <td><b>507 Plum / sand</b></td> <td>%.2f</td> <td>%s</td> </tr> 
			<tr> <td><b>508 Hog Plum Millet</b></td> <td>%.2f</td> <td>%s</td> </tr>  
			<tr> <td><b>509 Elaeocarpus serratus grain</b></td> <td>%.2f</td> <td>%s</td> </tr>  
			<tr> <td><b>510 Date corn</b></td> <td>%.2f</td> <td>%s</td> </tr>    
			<tr> <td><b>511 Palmyra</b></td> <td>%.2f</td> <td>%s</td> </tr>   
			<tr> <td><b>512 LentBellil</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>513 Limonia acidissima</b></td> <td>%.2f</td> <td>%s</td> </tr> 
			<tr> <td><b>514 Rose Apple</b></td> <td>%.2f</td> <td>%s</td> </tr>    
			<tr> <td><b>515 Carissa carandas</b></td> <td>%.2f</td> <td>%s</td> </tr>    
			<tr> <td><b>516 Sugar apple</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>517 custard apple</b></td> <td>%.2f</td> <td>%s</td> </tr>      
			<tr> <td><b>518 Pomegranate</b></td> <td>%.2f</td> <td>%s</td> </tr>       
			<tr> <td><b>519 Sapodilla</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>520 Monkey Jack</b></td> <td>%.2f</td> <td>%s</td> </tr>   
			<tr> <td><b>521 Averrhoa Carambola</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>522 Tamarind</b></td> <td>%.2f</td> <td>%s</td> </tr>       
			<tr> <td><b>523 Lemon</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>524 Grapefruit</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>525 Indian gooseberry</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>526 Baccaurea motleyana</b></td> <td>%.2f</td> <td>%s</td> </tr>     
			<tr> <td><b>527 Otaheite gooseberry</b></td> <td>%.2f</td> <td>%s</td> </tr>     
			<tr> <td><b>528 Elephant Apple</b></td> <td>%.2f</td> <td>%s</td> </tr> 
			<tr> <td><b>529 Orange</b></td> <td>%.2f</td> <td>%s</td> </tr>      
			<tr> <td><b>530 Citrus macroptera</b></td> <td>%.2f</td> <td>%s</td> </tr>
			<tr> <td><b>531 Citrus × sinensi</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>532 Bilimb</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>533 Velvet apple rice</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>534 Ficus</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>535 Dragon</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>536 Rambutan</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>537 Other fruits</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>538 Battle Leaf</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>539 Areca Catechu</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>540 Tea</b></td> <td>%.2f</td> <td>%s</td> </tr>                  
			<tr> <td><b>541 Other addicts</b></td> <td>%.2f</td> <td>%s</td> </tr>         
			<tr> <td><b>542 Cinnamomum tamala</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>543 Cinnamon</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>544 Cardamom</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>545 Other spices</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>546 Bamboo</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>547 Cane</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>548 Hardwood Tree</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>549 Rain Tree</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>550 Mahogany</b></td> <td>%.2f</td> <td>%s</td> </tr>      
			<tr> <td><b>551 Tectona Grandis</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>552 Auri</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>553 Gum trees</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>554 Baby tree</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>555 Gmelina arborea</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>556 Rhizophora apiculata</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>557 Banyan</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>559 Typha</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>560 Polyalthia longifolia</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>561 Alstonia scholaris</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>562 Jarul</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>563 Gum arabic</b></td> <td>%.2f</td> <td>%s</td> </tr>             
			<tr> <td><b>564 Monkey Jack</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>565 Golden Shower</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>566 Natural rubber</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>567 Aquilaria malaccensis</b></td> <td>%.2f</td> <td>%s</td> </tr>
			<tr> <td><b>568 Sakhua or shala</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>569 Artocarpus chama</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>570 Albizia lebbeck</b></td> <td>%.2f</td> <td>%s</td> </tr>                  
			<tr> <td><b>571 Other wood and forestry</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>572 Neem</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>573 Arjun tree</b></td> <td>%.2f</td> <td>%s</td> </tr>             
			<tr> <td><b>574 bahera or beleric</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>575 Myrobalan</b></td> <td>%.2f</td> <td>%s</td> </tr>       
			<tr> <td><b>577 Drumstick tree</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>579 Rose</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>580 Bakul</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>581 Neolamarckia cadamba</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>582 Royal poinciana</b></td> <td>%.2f</td> <td>%s</td> </tr>            
			<tr> <td><b>583 </b></td> <td>%.2f</td> <td>%s</td> </tr>            
			<tr> <td><b>584 Silk-cotton</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>585 Mulberry</b></td> <td>%.2f</td> <td>%s</td> </tr>                  
		</tbody>
	</table>
	</div>
	`,
		getTableGenerationName(q.TableNumber),
		p.Sprintf("%d", c.NumberOfFarmHoldings),
		p.Sprintf("%.2f", c.CropArea),
		c.p501a, c.PercentageOfPermantCropArea("P501"),
		c.p502a, c.PercentageOfPermantCropArea("P502"),
		c.p503a, c.PercentageOfPermantCropArea("P503"),
		c.p504a, c.PercentageOfPermantCropArea("P504"),
		c.p505a, c.PercentageOfPermantCropArea("P505"),
		c.p506a, c.PercentageOfPermantCropArea("P506"),
		c.p507a, c.PercentageOfPermantCropArea("P507"),
		c.p508a, c.PercentageOfCropArea("P508"),
		c.p509a, c.PercentageOfCropArea("P509"),
		c.p510a, c.PercentageOfCropArea("P510"),
		c.p511a, c.PercentageOfCropArea("P511"),
		c.p512a, c.PercentageOfCropArea("P512"),
		c.p513a, c.PercentageOfCropArea("P513"),
		c.p514a, c.PercentageOfCropArea("P514"),
		c.p515a, c.PercentageOfCropArea("P515"),
		c.p516a, c.PercentageOfCropArea("P516"),
		c.p517a, c.PercentageOfCropArea("P517"),
		c.p518a, c.PercentageOfCropArea("P518"),
		c.p519a, c.PercentageOfCropArea("P519"),
		c.p520a, c.PercentageOfCropArea("P520"),
		c.p521a, c.PercentageOfCropArea("P521"),
		c.p522a, c.PercentageOfCropArea("P522"),
		c.p523a, c.PercentageOfCropArea("P523"),
		c.p524a, c.PercentageOfCropArea("P524"),
		c.p525a, c.PercentageOfCropArea("P525"),
		c.p526a, c.PercentageOfCropArea("P526"),
		c.p527a, c.PercentageOfCropArea("P527"),
		c.p528a, c.PercentageOfCropArea("P528"),
		c.p529a, c.PercentageOfCropArea("P529"),
		c.p530a, c.PercentageOfCropArea("P530"),
		c.p531a, c.PercentageOfCropArea("P531"),
		c.p532a, c.PercentageOfCropArea("P532"),
		c.p533a, c.PercentageOfCropArea("P533"),
		c.p534a, c.PercentageOfCropArea("P534"),
		c.p535a, c.PercentageOfCropArea("P535"),
		c.p536a, c.PercentageOfCropArea("P536"),
		c.p537a, c.PercentageOfCropArea("P537"),
		c.p538a, c.PercentageOfCropArea("P538"),
		c.p539a, c.PercentageOfCropArea("P539"),
		c.p540a, c.PercentageOfCropArea("P540"),
		c.p541a, c.PercentageOfCropArea("P541"),
		c.p542a, c.PercentageOfCropArea("P542"),
		c.p543a, c.PercentageOfCropArea("P543"),
		c.p544a, c.PercentageOfCropArea("P544"),
		c.p545a, c.PercentageOfCropArea("P545"),
		c.p546a, c.PercentageOfCropArea("P546"),
		c.p547a, c.PercentageOfCropArea("P547"),
		c.p548a, c.PercentageOfCropArea("P548"),
		c.p549a, c.PercentageOfCropArea("P549"),
		c.p550a, c.PercentageOfCropArea("P550"),
		c.p551a, c.PercentageOfCropArea("P551"),
		c.p552a, c.PercentageOfCropArea("P552"),
		c.p553a, c.PercentageOfCropArea("P553"),
		c.p554a, c.PercentageOfCropArea("P554"),
		c.p555a, c.PercentageOfCropArea("P555"),
		c.p556a, c.PercentageOfCropArea("P556"),
		c.p557a, c.PercentageOfCropArea("P557"),
		c.p559a, c.PercentageOfCropArea("P559"),
		c.p560a, c.PercentageOfCropArea("P560"),
		c.p561a, c.PercentageOfCropArea("P561"),
		c.p562a, c.PercentageOfCropArea("P562"),
		c.p563a, c.PercentageOfCropArea("P563"),
		c.p564a, c.PercentageOfCropArea("P564"),
		c.p565a, c.PercentageOfCropArea("P565"),
		c.p566a, c.PercentageOfCropArea("P566"),
		c.p567a, c.PercentageOfCropArea("P567"),
		c.p568a, c.PercentageOfCropArea("P568"),
		c.p569a, c.PercentageOfCropArea("P569"),
		c.p570a, c.PercentageOfCropArea("P570"),
		c.p571a, c.PercentageOfCropArea("P571"),
		c.p572a, c.PercentageOfCropArea("P572"),
		c.p573a, c.PercentageOfCropArea("P573"),
		c.p574a, c.PercentageOfCropArea("P574"),
		c.p575a, c.PercentageOfCropArea("P575"),
		c.p577a, c.PercentageOfCropArea("P577"),
		c.p579a, c.PercentageOfCropArea("P579"),
		c.p580a, c.PercentageOfCropArea("P580"),
		c.p581a, c.PercentageOfCropArea("P581"),
		c.p582a, c.PercentageOfCropArea("P582"),
		c.p584a, c.PercentageOfCropArea("P584"),
		c.p585a, c.PercentageOfCropArea("P585"),
	)

	return
}
