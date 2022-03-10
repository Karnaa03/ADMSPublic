SELECT sum(t101) as t101,
    sum(t102) as t102,
    sum(sf + mf + lf) as holfings,
    sum(c13) as crop_area,
    sum(t101 + t102) as total_temporary_crop_area,
    sum(t101) / sum(t101 + t102) * 100 as percentage_t101
FROM agregateds
WHERE geocode ~ '20.46.43.142.*';
--@block
SELECT c13,
    case
        c13
        when 0 then 0,
        else sum(
            t101 + t102 + t103 + t104 + t105 + t112 + t113 + t114 + t121 + t122 + t123 + t124 + t125 + t127 + t128 + t129 + t130 + t131 + t132 + t134 + t135 + t157 + t158 + t159 + t160 + t161 + t167 + t169 + t175 + t176 + t177 + t179 + t182 + t185 + t106 + t107 + t108 + t109 + t110 + t111 + t115 + t116 + t117 + t118 + t119 + t120 + t126 + t133 + t136 + t137 + t138 + t139 + t140 + t141 + t142 + t143 + t144 + t145 + t146 + t147 + t148 + t149 + t150 + t151 + t152 + t153 + t154 + t155 + t156 + t162 + t163 + t164 + t165 + t166 + t168 + t170 + t171 + t172 + t173 + t174 + t178 + t180 + t181 + t183 + t184 + t186 + t187 + t188 + t189 + t190 + t191 + t192 + t193 + t194 + t195 + t196 + t197 + t198 + t199 + t200 + t201 + t202 + t203
        ) / c13
    END as c13,
    rmo
FROM agregateds
where subpath(geocode, 0, 2) = '20.46'
    AND true = true
GROUP BY rmo,
    c13;
--@block
select *
from agregateds
where c13 = 0;
--@block
SELECT sum(
        t101 + t102 + t103 + t104 + t105 + t112 + t113 + t114 + t121 + t122 + t123 + t124 + t125 + t127 + t128 + t129 + t130 + t131 + t132 + t134 + t135 + t157 + t158 + t159 + t160 + t161 + t167 + t169 + t175 + t176 + t177 + t179 + t182 + t185 + t106 + t107 + t108 + t109 + t110 + t111 + t115 + t116 + t117 + t118 + t119 + t120 + t126 + t133 + t136 + t137 + t138 + t139 + t140 + t141 + t142 + t143 + t144 + t145 + t146 + t147 + t148 + t149 + t150 + t151 + t152 + t153 + t154 + t155 + t156 + t162 + t163 + t164 + t165 + t166 + t168 + t170 + t171 + t172 + t173 + t174 + t178 + t180 + t181 + t183 + t184 + t186 + t187 + t188 + t189 + t190 + t191 + t192 + t193 + t194 + t195 + t196 + t197 + t198 + t199 + t200 + t201 + t202 + t203
    ) / c13 as data,
    rmo
FROM agregateds
where subpath(geocode, 0, 1) = '20'
    AND true = true
GROUP BY rmo,
    c13;