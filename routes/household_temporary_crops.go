package routes

import (
	"fmt"
)

func (srv *Server) FormatHouseholdTemporaryCrops(division, district, upazilla, union, mouza string, q *searchQuery, geoLocation string) (tableAndDonut string, err error) {
	c, err := srv.Db.GetTemporaryCrops(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	//p := message.NewPrinter(language.English)

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> </small></h4>
	
	
	<table id="datatable-buttons" class="table table-striped">
		<thead>
		<tr>
	<th class="text-wrap" style="width: 500px;">Data for table name : %s</th>
	<th></th>
	<th></th>
	</tr>
		
		</thead>
		<tbody>
		<tr>
		<th>Crop code and name</th>
		<th>Total temporary crop area (acres)</th>
		<th>Percentage of crop area (acres)</th>
	</tr>
			<tr> <td><b>101 Aush</b> <td>%s</td> <td>%s</td> </tr>        
			<tr> <td><b>102 Almond</b></td> <td>%s</td> <td>%s</td> </tr>        
			<tr> <td><b>103 Boro</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>104 Wheat</b></td> <td>%s</td> <td>%s</td> </tr>         
			<tr> <td><b>105 Maize</b></td> <td>%s</td> <td>%s</td> </tr>         
			<tr> <td><b>106 Foxtail millet</b></td> <td>%s</td> <td>%s</td> </tr>
			<tr> <td><b>107 Barley / sand</b></td> <td>%s</td> <td>%s</td> </tr> 
			<tr> <td><b>108 Proso Millet</b></td> <td>%s</td> <td>%s</td> </tr>  
			<tr> <td><b>109 Millet grain</b></td> <td>%s</td> <td>%s</td> </tr>  
			<tr> <td><b>110 Broom corn</b></td> <td>%s</td> <td>%s</td> </tr>    
			<tr> <td><b>111 Other grain</b></td> <td>%s</td> <td>%s</td> </tr>   
			<tr> <td><b>112 Lentil</b></td> <td>%s</td> <td>%s</td> </tr>        
			<tr> <td><b>113 Saffran Pulse</b></td> <td>%s</td> <td>%s</td> </tr> 
			<tr> <td><b>114 Moog Pulse</b></td> <td>%s</td> <td>%s</td> </tr>    
			<tr> <td><b>115 Black gram</b></td> <td>%s</td> <td>%s</td> </tr>    
			<tr> <td><b>116 Pea</b></td> <td>%s</td> <td>%s</td> </tr>           
			<tr> <td><b>117 Chickpea</b></td> <td>%s</td> <td>%s</td> </tr>      
			<tr> <td><b>118 Aarahar</b></td> <td>%s</td> <td>%s</td> </tr>       
			<tr> <td><b>119 Fallon</b></td> <td>%s</td> <td>%s</td> </tr>        
			<tr> <td><b>120 Other pulse</b></td> <td>%s</td> <td>%s</td> </tr>   
			<tr> <td><b>121 Potato</b></td> <td>%s</td> <td>%s</td> </tr>        
			<tr> <td><b>122 Brinjal</b></td> <td>%s</td> <td>%s</td> </tr>       
			<tr> <td><b>123 Radish</b></td> <td>%s</td> <td>%s</td> </tr>        
			<tr> <td><b>124 Bean</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>125 Tomato</b></td> <td>%s</td> <td>%s</td> </tr>        
			<tr> <td><b>126 Chichenga</b></td> <td>%s</td> <td>%s</td> </tr>     
			<tr> <td><b>127 Multitude</b></td> <td>%s</td> <td>%s</td> </tr>     
			<tr> <td><b>128 Ladies Finger</b></td> <td>%s</td> <td>%s</td> </tr> 
			<tr> <td><b>129 Cucumber</b></td> <td>%s</td> <td>%s</td> </tr>      
			<tr> <td><b>130 Bitter Gourd / Momordica / Charantia</b></td> <td>%s</td> <td>%s</td> </tr>
			<tr> <td><b>131 Gourd</b></td> <td>%s</td> <td>%s</td> </tr>                 
			<tr> <td><b>132 Pumpkin</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>133 Pumpkin rice</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>134 Cauliflower</b></td> <td>%s</td> <td>%s</td> </tr>           
			<tr> <td><b>135 Cabbage</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>136 Broccoli</b></td> <td>%s</td> <td>%s</td> </tr>              
			<tr> <td><b>137 Cucumber</b></td> <td>%s</td> <td>%s</td> </tr>              
			<tr> <td><b>138 Sweet potato</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>139 Stalk</b></td> <td>%s</td> <td>%s</td> </tr>                 
			<tr> <td><b>140 Taro</b></td> <td>%s</td> <td>%s</td> </tr>                  
			<tr> <td><b>141 Yardlong bean</b></td> <td>%s</td> <td>%s</td> </tr>         
			<tr> <td><b>142 Jhinga</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>143 Carrots</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>144 Kohlrabi</b></td> <td>%s</td> <td>%s</td> </tr>              
			<tr> <td><b>145 Turnip</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>146 Cumin</b></td> <td>%s</td> <td>%s</td> </tr>                 
			<tr> <td><b>147 Peppers</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>148 Sponge gourd</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>149 Beetroot</b></td> <td>%s</td> <td>%s</td> </tr>              
			<tr> <td><b>150 Other vegetables</b></td> <td>%s</td> <td>%s</td> </tr>      
			<tr> <td><b>151 Reddish</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>152 Indian spinach</b></td> <td>%s</td> <td>%s</td> </tr>        
			<tr> <td><b>153 Spinach</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>154 Mint leaves</b></td> <td>%s</td> <td>%s</td> </tr>           
			<tr> <td><b>155 lettuce leaf</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>156 Others leaf</b></td> <td>%s</td> <td>%s</td> </tr>           
			<tr> <td><b>157 Onion</b></td> <td>%s</td> <td>%s</td> </tr>                 
			<tr> <td><b>158 Garlic</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>159 Ginger</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>160 Turmeric</b></td> <td>%s</td> <td>%s</td> </tr>              
			<tr> <td><b>161 Chili</b></td> <td>%s</td> <td>%s</td> </tr>                 
			<tr> <td><b>162 Coriander</b></td> <td>%s</td> <td>%s</td> </tr>             
			<tr> <td><b>163 Black cumin</b></td> <td>%s</td> <td>%s</td> </tr>           
			<tr> <td><b>164 Fennel</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>165 Cumin</b></td> <td>%s</td> <td>%s</td> </tr>                 
			<tr> <td><b>166 Other spices  national</b></td> <td>%s</td> <td>%s</td> </tr>
			<tr> <td><b>167 Mustard</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>168 Soybean</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>169 Nuts</b></td> <td>%s</td> <td>%s</td> </tr>                  
			<tr> <td><b>170 Sesame</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>171 Linseed</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>172 sunflower</b></td> <td>%s</td> <td>%s</td> </tr>             
			<tr> <td><b>173 Castor</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>174 Other oil seeds</b></td> <td>%s</td> <td>%s</td> </tr>       
			<tr> <td><b>175 Banana</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>176 Papaya</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>177 Water Melon</b></td> <td>%s</td> <td>%s</td> </tr>           
			<tr> <td><b>178 Melons</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>179 Pine Apple</b></td> <td>%s</td> <td>%s</td> </tr>            
			<tr> <td><b>180 Strawberry</b></td> <td>%s</td> <td>%s</td> </tr>            
			<tr> <td><b>181 Other Fruits</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>182 Jute</b></td> <td>%s</td> <td>%s</td> </tr>                  
			<tr> <td><b>183 Cotton</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>184 Other fibers</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>185 Sugar Cane</b></td> <td>%s</td> <td>%s</td> </tr>            
			<tr> <td><b>186 Other Sugars</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>187 Tobacco</b></td> <td>%s</td> <td>%s</td> </tr>               
			<tr> <td><b>188 Other drugs</b></td> <td>%s</td> <td>%s</td> </tr>           
			<tr> <td><b>189 Aloe vera</b></td> <td>%s</td> <td>%s</td> </tr>             
			<tr> <td><b>190 Other medicinal "</b></td> <td>%s</td> <td>%s</td> </tr>     
			<tr> <td><b>191 Tuberose</b></td> <td>%s</td> <td>%s</td> </tr>              
			<tr> <td><b>192 Marigold</b></td> <td>%s</td> <td>%s</td> </tr>              
			<tr> <td><b>193 Chrysanthemum</b></td> <td>%s</td> <td>%s</td> </tr>         
			<tr> <td><b>194 Dahlia</b></td> <td>%s</td> <td>%s</td> </tr>                
			<tr> <td><b>195 Gladiolus</b></td> <td>%s</td> <td>%s</td> </tr>             
			<tr> <td><b>196 Transvaal daisy</b></td> <td>%s</td> <td>%s</td> </tr>       
			<tr> <td><b>197 Other flowers</b></td> <td>%s</td> <td>%s</td> </tr>         
			<tr> <td><b>198 Sun grass</b></td> <td>%s</td> <td>%s</td> </tr>             
			<tr> <td><b>199 Dhaincha</b></td> <td>%s</td> <td>%s</td> </tr>              
			<tr> <td><b>200 Other fuels</b></td> <td>%s</td> <td>%s</td> </tr>           
			<tr> <td><b>201 Napier grass</b></td> <td>%s</td> <td>%s</td> </tr>          
			<tr> <td><b>202 Other cow-Foods</b></td> <td>%s</td> <td>%s</td> </tr>       
			<tr> <td><b>203 Seeded</b></td> <td>%s</td> <td>%s</td> </tr>                
		</tbody>
	</table>
	</div>
	<h7>Source: Agriculture Census 2019, Bangladesh Bureau of Statistics.</h7>
	`,
		fmt.Sprintf("%s Geo CODE : %s", getTableGenerationName(q.TableNumber), geoLocation),

		FormatFloat(c.T101, 2), c.PercentageOfCropArea("T101"),
		FormatFloat(c.T102, 2), c.PercentageOfCropArea("T102"),
		FormatFloat(c.T103, 2), c.PercentageOfCropArea("T103"),
		FormatFloat(c.T104, 2), c.PercentageOfCropArea("T104"),
		FormatFloat(c.T105, 2), c.PercentageOfCropArea("T105"),
		FormatFloat(c.T106, 2), c.PercentageOfCropArea("T106"),
		FormatFloat(c.T107, 2), c.PercentageOfCropArea("T107"),
		FormatFloat(c.T108, 2), c.PercentageOfCropArea("T108"),
		FormatFloat(c.T109, 2), c.PercentageOfCropArea("T109"),
		FormatFloat(c.T110, 2), c.PercentageOfCropArea("T110"),
		FormatFloat(c.T111, 2), c.PercentageOfCropArea("T111"),
		FormatFloat(c.T112, 2), c.PercentageOfCropArea("T112"),
		FormatFloat(c.T113, 2), c.PercentageOfCropArea("T113"),
		FormatFloat(c.T114, 2), c.PercentageOfCropArea("T114"),
		FormatFloat(c.T115, 2), c.PercentageOfCropArea("T115"),
		FormatFloat(c.T116, 2), c.PercentageOfCropArea("T116"),
		FormatFloat(c.T117, 2), c.PercentageOfCropArea("T117"),
		FormatFloat(c.T118, 2), c.PercentageOfCropArea("T118"),
		FormatFloat(c.T119, 2), c.PercentageOfCropArea("T119"),
		FormatFloat(c.T120, 2), c.PercentageOfCropArea("T120"),
		FormatFloat(c.T121, 2), c.PercentageOfCropArea("T121"),
		FormatFloat(c.T122, 2), c.PercentageOfCropArea("T122"),
		FormatFloat(c.T123, 2), c.PercentageOfCropArea("T123"),
		FormatFloat(c.T124, 2), c.PercentageOfCropArea("T124"),
		FormatFloat(c.T125, 2), c.PercentageOfCropArea("T125"),
		FormatFloat(c.T126, 2), c.PercentageOfCropArea("T126"),
		FormatFloat(c.T127, 2), c.PercentageOfCropArea("T127"),
		FormatFloat(c.T128, 2), c.PercentageOfCropArea("T128"),
		FormatFloat(c.T129, 2), c.PercentageOfCropArea("T129"),
		FormatFloat(c.T130, 2), c.PercentageOfCropArea("T130"),
		FormatFloat(c.T131, 2), c.PercentageOfCropArea("T131"),
		FormatFloat(c.T132, 2), c.PercentageOfCropArea("T132"),
		FormatFloat(c.T133, 2), c.PercentageOfCropArea("T133"),
		FormatFloat(c.T134, 2), c.PercentageOfCropArea("T134"),
		FormatFloat(c.T135, 2), c.PercentageOfCropArea("T135"),
		FormatFloat(c.T136, 2), c.PercentageOfCropArea("T136"),
		FormatFloat(c.T137, 2), c.PercentageOfCropArea("T137"),
		FormatFloat(c.T138, 2), c.PercentageOfCropArea("T138"),
		FormatFloat(c.T139, 2), c.PercentageOfCropArea("T139"),
		FormatFloat(c.T140, 2), c.PercentageOfCropArea("T140"),
		FormatFloat(c.T141, 2), c.PercentageOfCropArea("T141"),
		FormatFloat(c.T142, 2), c.PercentageOfCropArea("T142"),
		FormatFloat(c.T143, 2), c.PercentageOfCropArea("T143"),
		FormatFloat(c.T144, 2), c.PercentageOfCropArea("T144"),
		FormatFloat(c.T145, 2), c.PercentageOfCropArea("T145"),
		FormatFloat(c.T146, 2), c.PercentageOfCropArea("T146"),
		FormatFloat(c.T147, 2), c.PercentageOfCropArea("T147"),
		FormatFloat(c.T148, 2), c.PercentageOfCropArea("T148"),
		FormatFloat(c.T149, 2), c.PercentageOfCropArea("T149"),
		FormatFloat(c.T150, 2), c.PercentageOfCropArea("T150"),
		FormatFloat(c.T151, 2), c.PercentageOfCropArea("T151"),
		FormatFloat(c.T152, 2), c.PercentageOfCropArea("T152"),
		FormatFloat(c.T153, 2), c.PercentageOfCropArea("T153"),
		FormatFloat(c.T154, 2), c.PercentageOfCropArea("T154"),
		FormatFloat(c.T155, 2), c.PercentageOfCropArea("T155"),
		FormatFloat(c.T156, 2), c.PercentageOfCropArea("T156"),
		FormatFloat(c.T157, 2), c.PercentageOfCropArea("T157"),
		FormatFloat(c.T158, 2), c.PercentageOfCropArea("T158"),
		FormatFloat(c.T159, 2), c.PercentageOfCropArea("T159"),
		FormatFloat(c.T160, 2), c.PercentageOfCropArea("T160"),
		FormatFloat(c.T161, 2), c.PercentageOfCropArea("T161"),
		FormatFloat(c.T162, 2), c.PercentageOfCropArea("T162"),
		FormatFloat(c.T163, 2), c.PercentageOfCropArea("T163"),
		FormatFloat(c.T164, 2), c.PercentageOfCropArea("T164"),
		FormatFloat(c.T165, 2), c.PercentageOfCropArea("T165"),
		FormatFloat(c.T166, 2), c.PercentageOfCropArea("T166"),
		FormatFloat(c.T167, 2), c.PercentageOfCropArea("T167"),
		FormatFloat(c.T168, 2), c.PercentageOfCropArea("T168"),
		FormatFloat(c.T169, 2), c.PercentageOfCropArea("T169"),
		FormatFloat(c.T170, 2), c.PercentageOfCropArea("T170"),
		FormatFloat(c.T171, 2), c.PercentageOfCropArea("T171"),
		FormatFloat(c.T172, 2), c.PercentageOfCropArea("T172"),
		FormatFloat(c.T173, 2), c.PercentageOfCropArea("T173"),
		FormatFloat(c.T174, 2), c.PercentageOfCropArea("T174"),
		FormatFloat(c.T175, 2), c.PercentageOfCropArea("T175"),
		FormatFloat(c.T176, 2), c.PercentageOfCropArea("T176"),
		FormatFloat(c.T177, 2), c.PercentageOfCropArea("T177"),
		FormatFloat(c.T178, 2), c.PercentageOfCropArea("T178"),
		FormatFloat(c.T179, 2), c.PercentageOfCropArea("T179"),
		FormatFloat(c.T180, 2), c.PercentageOfCropArea("T180"),
		FormatFloat(c.T181, 2), c.PercentageOfCropArea("T181"),
		FormatFloat(c.T182, 2), c.PercentageOfCropArea("T182"),
		FormatFloat(c.T183, 2), c.PercentageOfCropArea("T183"),
		FormatFloat(c.T184, 2), c.PercentageOfCropArea("T184"),
		FormatFloat(c.T185, 2), c.PercentageOfCropArea("T185"),
		FormatFloat(c.T186, 2), c.PercentageOfCropArea("T186"),
		FormatFloat(c.T187, 2), c.PercentageOfCropArea("T187"),
		FormatFloat(c.T188, 2), c.PercentageOfCropArea("T188"),
		FormatFloat(c.T189, 2), c.PercentageOfCropArea("T189"),
		FormatFloat(c.T190, 2), c.PercentageOfCropArea("T190"),
		FormatFloat(c.T191, 2), c.PercentageOfCropArea("T191"),
		FormatFloat(c.T192, 2), c.PercentageOfCropArea("T192"),
		FormatFloat(c.T193, 2), c.PercentageOfCropArea("T193"),
		FormatFloat(c.T194, 2), c.PercentageOfCropArea("T194"),
		FormatFloat(c.T195, 2), c.PercentageOfCropArea("T195"),
		FormatFloat(c.T196, 2), c.PercentageOfCropArea("T196"),
		FormatFloat(c.T197, 2), c.PercentageOfCropArea("T197"),
		FormatFloat(c.T198, 2), c.PercentageOfCropArea("T198"),
		FormatFloat(c.T199, 2), c.PercentageOfCropArea("T199"),
		FormatFloat(c.T200, 2), c.PercentageOfCropArea("T200"),
		FormatFloat(c.T201, 2), c.PercentageOfCropArea("T201"),
		FormatFloat(c.T202, 2), c.PercentageOfCropArea("T202"),
		FormatFloat(c.T203, 2), c.PercentageOfCropArea("T203"),
	)

	return
}
