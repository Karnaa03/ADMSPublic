	SELECT sum(edu) as no_education,
    	(sum(edu1) + sum(edu2) + sum(edu3) + sum(edu4) + sum(edu5)) as class_I_V,
    	(sum(edu6) + sum(edu7) + sum(edu8) + sum(edu9)) as class_VI_IX,
    	sum(edu10) as Scc_Passed,
    	sum(edu12) as Hsc_Passed,
    	sum(edu15) as Degree_Passed,
    	sum(edu18) as Master_Passed,
    	(
    	    sum(edu) + sum(edu1) + sum(edu2) + sum(edu3) + sum(edu4) + sum(edu5) + sum(edu6) + sum(edu7) + sum(edu8) + sum(edu9) + sum(edu10) + sum(edu12) + sum(edu15) + sum(edu18)
    	) as Total_Education,
    	sum(occ) as Agriculture,
    	sum(occ2) as Industry,
    	sum(occ3) as Service,
    	sum(occ4) as Business,
    	sum(occ5) as Other,
    	(sum(occ) + sum(occ2) + sum(occ3) + sum(occ4) + sum(occ5)) as Total_Occupation,
    	sum(c01m) as Household_Member_Male,
    	sum(c01f) as Household_Member_Female,
    	sum(c01h) as Household_Member_Hijra,
    	(sum(c01m) + sum(c01f) + sum(c01h)) as Household_Member_Total,
    	(sum(c02m) + sum(c03m)) as Household_Worker_Male,
    	(sum(c02f) + sum(c03f)) as Household_Worker_Female,
    	(sum(c02h) + sum(c03h)) as Household_Worker_Hijra,
    	(sum(c02m) + sum(c03m) + sum(c02f) + sum(c03f) + sum(c02h) + sum(c03h)) as Household_Worker_Total,
    	sum(c02m) as Household_Worker_10_14_Male,
    	sum(c02f) as Household_Worker_10_14_Female,
    	sum(c02h) as Household_Worker_10_14_Hijra,
    	(sum(c02m) + sum(c02f) + sum(c02h)) as Household_Worker_10_14_Total,
    	sum(c03m) as Household_Worker_15_Plus_Male,
    	sum(c03f) as Household_Worker_15_Plus_Female,
    	sum(c03h) as Household_Worker_15_Plus_Hijra,
    	(sum(c03m) + sum(c03f) + sum(c03h)) as Household_Worker_15_Plus_Total
	FROM aggregates
	WHERE hh_a = 1 AND subpath(geocode, 0, 0) = '';

--@block
SELECT sum(hh_sno) as data FROM aggregates where subpath(geocode, 0,0) = '' ;	


--@block
Select  SUM(c02m+c02f+c02h+c03m+c03f+c03h)::NUMERIC as somme, 
SUM(c02m+c02f+c02h+c03m+c03f+c03h)::NUMERIC/SUM(hh_sno)::NUMERIC as ratio, 
SUM(hh_sno)::NUMERIC as somme_hno, rmo
from aggregates
where true = true GROUP BY rmo;