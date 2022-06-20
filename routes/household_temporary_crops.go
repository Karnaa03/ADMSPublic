package routes

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (srv *Server) FormatHouseholdTemporaryCrops(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	c, err := srv.Db.GetTemporaryCrops(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	p := message.NewPrinter(language.English)

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> </small></h4>
	<h5>Data for table name : %s</h5>
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
			<tr> <td><b>101 Aush</b> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>102 Almond</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>103 Boro</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>104 Wheat</b></td> <td>%.2f</td> <td>%s</td> </tr>         
			<tr> <td><b>105 Maize</b></td> <td>%.2f</td> <td>%s</td> </tr>         
			<tr> <td><b>106 Foxtail millet</b></td> <td>%.2f</td> <td>%s</td> </tr>
			<tr> <td><b>107 Barley / sand</b></td> <td>%.2f</td> <td>%s</td> </tr> 
			<tr> <td><b>108 Proso Millet</b></td> <td>%.2f</td> <td>%s</td> </tr>  
			<tr> <td><b>109 Millet grain</b></td> <td>%.2f</td> <td>%s</td> </tr>  
			<tr> <td><b>110 Broom corn</b></td> <td>%.2f</td> <td>%s</td> </tr>    
			<tr> <td><b>111 Other grain</b></td> <td>%.2f</td> <td>%s</td> </tr>   
			<tr> <td><b>112 Lentil</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>113 Saffran Pulse</b></td> <td>%.2f</td> <td>%s</td> </tr> 
			<tr> <td><b>114 Moog Pulse</b></td> <td>%.2f</td> <td>%s</td> </tr>    
			<tr> <td><b>115 Black gram</b></td> <td>%.2f</td> <td>%s</td> </tr>    
			<tr> <td><b>116 Pea</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>117 Chickpea</b></td> <td>%.2f</td> <td>%s</td> </tr>      
			<tr> <td><b>118 Aarahar</b></td> <td>%.2f</td> <td>%s</td> </tr>       
			<tr> <td><b>119 Fallon</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>120 Other pulse</b></td> <td>%.2f</td> <td>%s</td> </tr>   
			<tr> <td><b>121 Potato</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>122 Brinjal</b></td> <td>%.2f</td> <td>%s</td> </tr>       
			<tr> <td><b>123 Radish</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>124 Bean</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>125 Tomato</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>126 Chichenga</b></td> <td>%.2f</td> <td>%s</td> </tr>     
			<tr> <td><b>127 Multitude</b></td> <td>%.2f</td> <td>%s</td> </tr>     
			<tr> <td><b>128 Ladies Finger</b></td> <td>%.2f</td> <td>%s</td> </tr> 
			<tr> <td><b>129 Cucumber</b></td> <td>%.2f</td> <td>%s</td> </tr>      
			<tr> <td><b>130 Bitter Gourd / Momordica / Charantia</b></td> <td>%.2f</td> <td>%s</td> </tr>
			<tr> <td><b>131 Gourd</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>132 Pumpkin</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>133 Pumpkin rice</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>134 Cauliflower</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>135 Cabbage</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>136 Broccoli</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>137 Cucumber</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>138 Sweet potato</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>139 Stalk</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>140 Taro</b></td> <td>%.2f</td> <td>%s</td> </tr>                  
			<tr> <td><b>141 Yardlong bean</b></td> <td>%.2f</td> <td>%s</td> </tr>         
			<tr> <td><b>142 Jhinga</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>143 Carrots</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>144 Kohlrabi</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>145 Turnip</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>146 Cumin</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>147 Peppers</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>148 Sponge gourd</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>149 Beetroot</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>150 Other vegetables</b></td> <td>%.2f</td> <td>%s</td> </tr>      
			<tr> <td><b>151 Reddish</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>152 Indian spinach</b></td> <td>%.2f</td> <td>%s</td> </tr>        
			<tr> <td><b>153 Spinach</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>154 Mint leaves</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>155 lettuce leaf</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>156 Others leaf</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>157 Onion</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>158 Garlic</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>159 Ginger</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>160 Turmeric</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>161 Chili</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>162 Coriander</b></td> <td>%.2f</td> <td>%s</td> </tr>             
			<tr> <td><b>163 Black cumin</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>164 Fennel</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>165 Cumin</b></td> <td>%.2f</td> <td>%s</td> </tr>                 
			<tr> <td><b>166 Other spices  national</b></td> <td>%.2f</td> <td>%s</td> </tr>
			<tr> <td><b>167 Mustard</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>168 Soybean</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>169 Nuts</b></td> <td>%.2f</td> <td>%s</td> </tr>                  
			<tr> <td><b>170 Sesame</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>171 Linseed</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>172 sunflower</b></td> <td>%.2f</td> <td>%s</td> </tr>             
			<tr> <td><b>173 Castor</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>174 Other oil seeds</b></td> <td>%.2f</td> <td>%s</td> </tr>       
			<tr> <td><b>175 Banana</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>176 Papaya</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>177 Water Melon</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>178 Melons</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>179 Pine Apple</b></td> <td>%.2f</td> <td>%s</td> </tr>            
			<tr> <td><b>180 Strawberry</b></td> <td>%.2f</td> <td>%s</td> </tr>            
			<tr> <td><b>181 Other Fruits</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>182 Jute</b></td> <td>%.2f</td> <td>%s</td> </tr>                  
			<tr> <td><b>183 Cotton</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>184 Other fibers</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>185 Sugar Cane</b></td> <td>%.2f</td> <td>%s</td> </tr>            
			<tr> <td><b>186 Other Sugars</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>187 Tobacco</b></td> <td>%.2f</td> <td>%s</td> </tr>               
			<tr> <td><b>188 Other drugs</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>189 Aloe vera</b></td> <td>%.2f</td> <td>%s</td> </tr>             
			<tr> <td><b>190 Other medicinal "</b></td> <td>%.2f</td> <td>%s</td> </tr>     
			<tr> <td><b>191 Tuberose</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>192 Marigold</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>193 Chrysanthemum</b></td> <td>%.2f</td> <td>%s</td> </tr>         
			<tr> <td><b>194 Dahlia</b></td> <td>%.2f</td> <td>%s</td> </tr>                
			<tr> <td><b>195 Gladiolus</b></td> <td>%.2f</td> <td>%s</td> </tr>             
			<tr> <td><b>196 Transvaal daisy</b></td> <td>%.2f</td> <td>%s</td> </tr>       
			<tr> <td><b>197 Other flowers</b></td> <td>%.2f</td> <td>%s</td> </tr>         
			<tr> <td><b>198 Sun grass</b></td> <td>%.2f</td> <td>%s</td> </tr>             
			<tr> <td><b>199 Dhaincha</b></td> <td>%.2f</td> <td>%s</td> </tr>              
			<tr> <td><b>200 Other fuels</b></td> <td>%.2f</td> <td>%s</td> </tr>           
			<tr> <td><b>201 Napier grass</b></td> <td>%.2f</td> <td>%s</td> </tr>          
			<tr> <td><b>202 Other cow-Foods</b></td> <td>%.2f</td> <td>%s</td> </tr>       
			<tr> <td><b>203 Seeded</b></td> <td>%.2f</td> <td>%s</td> </tr>                
		</tbody>
	</table>
	</div>
	`,
		getTableGenerationName(q.TableNumber),
		p.Sprintf("%d", c.NumberOfFarmHoldings),
		p.Sprintf("%.2f", c.CropArea),
		c.T101, c.PercentageOfCropArea("T101"),
		c.T102, c.PercentageOfCropArea("T102"),
		c.T103, c.PercentageOfCropArea("T103"),
		c.T104, c.PercentageOfCropArea("T104"),
		c.T105, c.PercentageOfCropArea("T105"),
		c.T106, c.PercentageOfCropArea("T106"),
		c.T107, c.PercentageOfCropArea("T107"),
		c.T108, c.PercentageOfCropArea("T108"),
		c.T109, c.PercentageOfCropArea("T109"),
		c.T110, c.PercentageOfCropArea("T110"),
		c.T111, c.PercentageOfCropArea("T111"),
		c.T112, c.PercentageOfCropArea("T112"),
		c.T113, c.PercentageOfCropArea("T113"),
		c.T114, c.PercentageOfCropArea("T114"),
		c.T115, c.PercentageOfCropArea("T115"),
		c.T116, c.PercentageOfCropArea("T116"),
		c.T117, c.PercentageOfCropArea("T117"),
		c.T118, c.PercentageOfCropArea("T118"),
		c.T119, c.PercentageOfCropArea("T119"),
		c.T120, c.PercentageOfCropArea("T120"),
		c.T121, c.PercentageOfCropArea("T121"),
		c.T122, c.PercentageOfCropArea("T122"),
		c.T123, c.PercentageOfCropArea("T123"),
		c.T124, c.PercentageOfCropArea("T124"),
		c.T125, c.PercentageOfCropArea("T125"),
		c.T126, c.PercentageOfCropArea("T126"),
		c.T127, c.PercentageOfCropArea("T127"),
		c.T128, c.PercentageOfCropArea("T128"),
		c.T129, c.PercentageOfCropArea("T129"),
		c.T130, c.PercentageOfCropArea("T130"),
		c.T131, c.PercentageOfCropArea("T131"),
		c.T132, c.PercentageOfCropArea("T132"),
		c.T133, c.PercentageOfCropArea("T133"),
		c.T134, c.PercentageOfCropArea("T134"),
		c.T135, c.PercentageOfCropArea("T135"),
		c.T136, c.PercentageOfCropArea("T136"),
		c.T137, c.PercentageOfCropArea("T137"),
		c.T138, c.PercentageOfCropArea("T138"),
		c.T139, c.PercentageOfCropArea("T139"),
		c.T140, c.PercentageOfCropArea("T140"),
		c.T141, c.PercentageOfCropArea("T141"),
		c.T142, c.PercentageOfCropArea("T142"),
		c.T143, c.PercentageOfCropArea("T143"),
		c.T144, c.PercentageOfCropArea("T144"),
		c.T145, c.PercentageOfCropArea("T145"),
		c.T146, c.PercentageOfCropArea("T146"),
		c.T147, c.PercentageOfCropArea("T147"),
		c.T148, c.PercentageOfCropArea("T148"),
		c.T149, c.PercentageOfCropArea("T149"),
		c.T150, c.PercentageOfCropArea("T150"),
		c.T151, c.PercentageOfCropArea("T151"),
		c.T152, c.PercentageOfCropArea("T152"),
		c.T153, c.PercentageOfCropArea("T153"),
		c.T154, c.PercentageOfCropArea("T154"),
		c.T155, c.PercentageOfCropArea("T155"),
		c.T156, c.PercentageOfCropArea("T156"),
		c.T157, c.PercentageOfCropArea("T157"),
		c.T158, c.PercentageOfCropArea("T158"),
		c.T159, c.PercentageOfCropArea("T159"),
		c.T160, c.PercentageOfCropArea("T160"),
		c.T161, c.PercentageOfCropArea("T161"),
		c.T162, c.PercentageOfCropArea("T162"),
		c.T163, c.PercentageOfCropArea("T163"),
		c.T164, c.PercentageOfCropArea("T164"),
		c.T165, c.PercentageOfCropArea("T165"),
		c.T166, c.PercentageOfCropArea("T166"),
		c.T167, c.PercentageOfCropArea("T167"),
		c.T168, c.PercentageOfCropArea("T168"),
		c.T169, c.PercentageOfCropArea("T169"),
		c.T170, c.PercentageOfCropArea("T170"),
		c.T171, c.PercentageOfCropArea("T171"),
		c.T172, c.PercentageOfCropArea("T172"),
		c.T173, c.PercentageOfCropArea("T173"),
		c.T174, c.PercentageOfCropArea("T174"),
		c.T175, c.PercentageOfCropArea("T175"),
		c.T176, c.PercentageOfCropArea("T176"),
		c.T177, c.PercentageOfCropArea("T177"),
		c.T178, c.PercentageOfCropArea("T178"),
		c.T179, c.PercentageOfCropArea("T179"),
		c.T180, c.PercentageOfCropArea("T180"),
		c.T181, c.PercentageOfCropArea("T181"),
		c.T182, c.PercentageOfCropArea("T182"),
		c.T183, c.PercentageOfCropArea("T183"),
		c.T184, c.PercentageOfCropArea("T184"),
		c.T185, c.PercentageOfCropArea("T185"),
		c.T186, c.PercentageOfCropArea("T186"),
		c.T187, c.PercentageOfCropArea("T187"),
		c.T188, c.PercentageOfCropArea("T188"),
		c.T189, c.PercentageOfCropArea("T189"),
		c.T190, c.PercentageOfCropArea("T190"),
		c.T191, c.PercentageOfCropArea("T191"),
		c.T192, c.PercentageOfCropArea("T192"),
		c.T193, c.PercentageOfCropArea("T193"),
		c.T194, c.PercentageOfCropArea("T194"),
		c.T195, c.PercentageOfCropArea("T195"),
		c.T196, c.PercentageOfCropArea("T196"),
		c.T197, c.PercentageOfCropArea("T197"),
		c.T198, c.PercentageOfCropArea("T198"),
		c.T199, c.PercentageOfCropArea("T199"),
		c.T200, c.PercentageOfCropArea("T200"),
		c.T201, c.PercentageOfCropArea("T201"),
		c.T202, c.PercentageOfCropArea("T202"),
		c.T203, c.PercentageOfCropArea("T203"),
	)

	return
}
