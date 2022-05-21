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
	<h7>Source: Bangladesh Bureau of Statistics. Report produced by Agriculture (Crops, Fisheries and Livestock) Census 2018 Project.</h7>
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
		c.P501a, c.PercentageOfPermantCropArea("P501a"),
		c.P502a, c.PercentageOfPermantCropArea("P502a"),
		c.P503a, c.PercentageOfPermantCropArea("P503a"),
		c.P504a, c.PercentageOfPermantCropArea("P504a"),
		c.P505a, c.PercentageOfPermantCropArea("P505a"),
		c.P506a, c.PercentageOfPermantCropArea("P506a"),
		c.P507a, c.PercentageOfPermantCropArea("P507a"),
		c.P508a, c.PercentageOfPermantCropArea("P508a"),
		c.P509a, c.PercentageOfPermantCropArea("P509a"),
		c.P510a, c.PercentageOfPermantCropArea("P510a"),
		c.P511a, c.PercentageOfPermantCropArea("P511a"),
		c.P512a, c.PercentageOfPermantCropArea("P512a"),
		c.P513a, c.PercentageOfPermantCropArea("P513a"),
		c.P514a, c.PercentageOfPermantCropArea("P514a"),
		c.P515a, c.PercentageOfPermantCropArea("P515a"),
		c.P516a, c.PercentageOfPermantCropArea("P516a"),
		c.P517a, c.PercentageOfPermantCropArea("P517a"),
		c.P518a, c.PercentageOfPermantCropArea("P518a"),
		c.P519a, c.PercentageOfPermantCropArea("P519a"),
		c.P520a, c.PercentageOfPermantCropArea("P520a"),
		c.P521a, c.PercentageOfPermantCropArea("P521a"),
		c.P522a, c.PercentageOfPermantCropArea("P522a"),
		c.P523a, c.PercentageOfPermantCropArea("P523a"),
		c.P524a, c.PercentageOfPermantCropArea("P524a"),
		c.P525a, c.PercentageOfPermantCropArea("P525a"),
		c.P526a, c.PercentageOfPermantCropArea("P526a"),
		c.P527a, c.PercentageOfPermantCropArea("P527a"),
		c.P528a, c.PercentageOfPermantCropArea("P528a"),
		c.P529a, c.PercentageOfPermantCropArea("P529a"),
		c.P530a, c.PercentageOfPermantCropArea("P530a"),
		c.P531a, c.PercentageOfPermantCropArea("P531a"),
		c.P532a, c.PercentageOfPermantCropArea("P532a"),
		c.P533a, c.PercentageOfPermantCropArea("P533a"),
		c.P534a, c.PercentageOfPermantCropArea("P534a"),
		c.P535a, c.PercentageOfPermantCropArea("P535a"),
		c.P536a, c.PercentageOfPermantCropArea("P536a"),
		c.P537a, c.PercentageOfPermantCropArea("P537a"),
		c.P538a, c.PercentageOfPermantCropArea("P538a"),
		c.P539a, c.PercentageOfPermantCropArea("P539a"),
		c.P540a, c.PercentageOfPermantCropArea("P540a"),
		c.P541a, c.PercentageOfPermantCropArea("P541a"),
		c.P542a, c.PercentageOfPermantCropArea("P542a"),
		c.P543a, c.PercentageOfPermantCropArea("P543a"),
		c.P544a, c.PercentageOfPermantCropArea("P544a"),
		c.P545a, c.PercentageOfPermantCropArea("P545a"),
		c.P546a, c.PercentageOfPermantCropArea("P546a"),
		c.P547a, c.PercentageOfPermantCropArea("P547a"),
		c.P548a, c.PercentageOfPermantCropArea("P548a"),
		c.P549a, c.PercentageOfPermantCropArea("P549a"),
		c.P550a, c.PercentageOfPermantCropArea("P550a"),
		c.P551a, c.PercentageOfPermantCropArea("P551a"),
		c.P552a, c.PercentageOfPermantCropArea("P552a"),
		c.P553a, c.PercentageOfPermantCropArea("P553a"),
		c.P554a, c.PercentageOfPermantCropArea("P554a"),
		c.P555a, c.PercentageOfPermantCropArea("P555a"),
		c.P556a, c.PercentageOfPermantCropArea("P556a"),
		c.P557a, c.PercentageOfPermantCropArea("P557a"),
		c.P559a, c.PercentageOfPermantCropArea("P559a"),
		c.P560a, c.PercentageOfPermantCropArea("P560a"),
		c.P561a, c.PercentageOfPermantCropArea("P561a"),
		c.P562a, c.PercentageOfPermantCropArea("P562a"),
		c.P563a, c.PercentageOfPermantCropArea("P563a"),
		c.P564a, c.PercentageOfPermantCropArea("P564a"),
		c.P565a, c.PercentageOfPermantCropArea("P565a"),
		c.P566a, c.PercentageOfPermantCropArea("P566a"),
		c.P567a, c.PercentageOfPermantCropArea("P567a"),
		c.P568a, c.PercentageOfPermantCropArea("P568a"),
		c.P569a, c.PercentageOfPermantCropArea("P569a"),
		c.P570a, c.PercentageOfPermantCropArea("P570a"),
		c.P571a, c.PercentageOfPermantCropArea("P571a"),
		c.P572a, c.PercentageOfPermantCropArea("P572a"),
		c.P573a, c.PercentageOfPermantCropArea("P573a"),
		c.P574a, c.PercentageOfPermantCropArea("P574a"),
		c.P575a, c.PercentageOfPermantCropArea("P575a"),
		c.P577a, c.PercentageOfPermantCropArea("P577a"),
		c.P579a, c.PercentageOfPermantCropArea("P579a"),
		c.P580a, c.PercentageOfPermantCropArea("P580a"),
		c.P581a, c.PercentageOfPermantCropArea("P581a"),
		c.P582a, c.PercentageOfPermantCropArea("P582a"),
		c.P584a, c.PercentageOfPermantCropArea("P584a"),
		c.P585a, c.PercentageOfPermantCropArea("P585a"),
	)

	return
}
