package routes

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (srv *Server) FormatHouseholdPermanentCrops(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	c, err := srv.Db.GetPermanantCrops(division, district, upazilla, union, mouza)
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
			<tr> <td><b>584 Silk-cotton</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>585 Mulberry</b></td> <td>%.2f</td> <td>%s</td> </tr>                  
		</tbody>
	</table>
	</div>
	`,
		getTableGenerationName(q.TableNumber),
		p.Sprintf("%d", c.NumberOfFarmHoldings),
		p.Sprintf("%.2f", c.CropArea),
		c.P501a, c.PercentageOfPermantCropArea("P501"),
		c.P502a, c.PercentageOfPermantCropArea("P502"),
		c.P503a, c.PercentageOfPermantCropArea("P503"),
		c.P504a, c.PercentageOfPermantCropArea("P504"),
		c.P505a, c.PercentageOfPermantCropArea("P505"),
		c.P506a, c.PercentageOfPermantCropArea("P506"),
		c.P507a, c.PercentageOfPermantCropArea("P507"),
		c.P508a, c.PercentageOfPermantCropArea("P508"),
		c.P509a, c.PercentageOfPermantCropArea("P509"),
		c.P510a, c.PercentageOfPermantCropArea("P510"),
		c.P511a, c.PercentageOfPermantCropArea("P511"),
		c.P512a, c.PercentageOfPermantCropArea("P512"),
		c.P513a, c.PercentageOfPermantCropArea("P513"),
		c.P514a, c.PercentageOfPermantCropArea("P514"),
		c.P515a, c.PercentageOfPermantCropArea("P515"),
		c.P516a, c.PercentageOfPermantCropArea("P516"),
		c.P517a, c.PercentageOfPermantCropArea("P517"),
		c.P518a, c.PercentageOfPermantCropArea("P518"),
		c.P519a, c.PercentageOfPermantCropArea("P519"),
		c.P520a, c.PercentageOfPermantCropArea("P520"),
		c.P521a, c.PercentageOfPermantCropArea("P521"),
		c.P522a, c.PercentageOfPermantCropArea("P522"),
		c.P523a, c.PercentageOfPermantCropArea("P523"),
		c.P524a, c.PercentageOfPermantCropArea("P524"),
		c.P525a, c.PercentageOfPermantCropArea("P525"),
		c.P526a, c.PercentageOfPermantCropArea("P526"),
		c.P527a, c.PercentageOfPermantCropArea("P527"),
		c.P528a, c.PercentageOfPermantCropArea("P528"),
		c.P529a, c.PercentageOfPermantCropArea("P529"),
		c.P530a, c.PercentageOfPermantCropArea("P530"),
		c.P531a, c.PercentageOfPermantCropArea("P531"),
		c.P532a, c.PercentageOfPermantCropArea("P532"),
		c.P533a, c.PercentageOfPermantCropArea("P533"),
		c.P534a, c.PercentageOfPermantCropArea("P534"),
		c.P535a, c.PercentageOfPermantCropArea("P535"),
		c.P536a, c.PercentageOfPermantCropArea("P536"),
		c.P537a, c.PercentageOfPermantCropArea("P537"),
		c.P538a, c.PercentageOfPermantCropArea("P538"),
		c.P539a, c.PercentageOfPermantCropArea("P539"),
		c.P540a, c.PercentageOfPermantCropArea("P540"),
		c.P541a, c.PercentageOfPermantCropArea("P541"),
		c.P542a, c.PercentageOfPermantCropArea("P542"),
		c.P543a, c.PercentageOfPermantCropArea("P543"),
		c.P544a, c.PercentageOfPermantCropArea("P544"),
		c.P545a, c.PercentageOfPermantCropArea("P545"),
		c.P546a, c.PercentageOfPermantCropArea("P546"),
		c.P547a, c.PercentageOfPermantCropArea("P547"),
		c.P548a, c.PercentageOfPermantCropArea("P548"),
		c.P549a, c.PercentageOfPermantCropArea("P549"),
		c.P550a, c.PercentageOfPermantCropArea("P550"),
		c.P551a, c.PercentageOfPermantCropArea("P551"),
		c.P552a, c.PercentageOfPermantCropArea("P552"),
		c.P553a, c.PercentageOfPermantCropArea("P553"),
		c.P554a, c.PercentageOfPermantCropArea("P554"),
		c.P555a, c.PercentageOfPermantCropArea("P555"),
		c.P556a, c.PercentageOfPermantCropArea("P556"),
		c.P557a, c.PercentageOfPermantCropArea("P557"),
		c.P559a, c.PercentageOfPermantCropArea("P559"),
		c.P560a, c.PercentageOfPermantCropArea("P560"),
		c.P561a, c.PercentageOfPermantCropArea("P561"),
		c.P562a, c.PercentageOfPermantCropArea("P562"),
		c.P563a, c.PercentageOfPermantCropArea("P563"),
		c.P564a, c.PercentageOfPermantCropArea("P564"),
		c.P565a, c.PercentageOfPermantCropArea("P565"),
		c.P566a, c.PercentageOfPermantCropArea("P566"),
		c.P567a, c.PercentageOfPermantCropArea("P567"),
		c.P568a, c.PercentageOfPermantCropArea("P568"),
		c.P569a, c.PercentageOfPermantCropArea("P569"),
		c.P570a, c.PercentageOfPermantCropArea("P570"),
		c.P571a, c.PercentageOfPermantCropArea("P571"),
		c.P572a, c.PercentageOfPermantCropArea("P572"),
		c.P573a, c.PercentageOfPermantCropArea("P573"),
		c.P574a, c.PercentageOfPermantCropArea("P574"),
		c.P575a, c.PercentageOfPermantCropArea("P575"),
		c.P577a, c.PercentageOfPermantCropArea("P577"),
		c.P579a, c.PercentageOfPermantCropArea("P579"),
		c.P580a, c.PercentageOfPermantCropArea("P580"),
		c.P581a, c.PercentageOfPermantCropArea("P581"),
		c.P582a, c.PercentageOfPermantCropArea("P582"),
		c.P584a, c.PercentageOfPermantCropArea("P584"),
		c.P585a, c.PercentageOfPermantCropArea("P585"),
	)

	return
}
