select edu,
    edu1 as NoEducation,
    edu2 as Class1,
    edu3 as Class2,
    edu3 as Class3,
    edu4 as Class4,
    edu5 as Class5,
    edu6 as Class6,
    edu7 as Class7,
    edu8 as Class8,
    edu9 as Class9,
    edu10 as SCC,
    edu12 as HSC,
    edu15 as BachelorEquivalent,
    edu18 as MastersEquivalentOrHigher,
    (
        edu1 + edu2 + edu3 + edu4 + edu5 + edu6 + edu7 + edu8 + edu9 + edu10 + edu12 + edu15 + edu18
    ) as Total
from agregateds;
-- @block test
select sum(occ) as occ,
    sum(occ2) as occ2,
    sum(occ3) as occ3,
    sum(occ4) as occ4,
    sum(occ5) as occ5,
    (
        sum(occ) + sum(occ2) + sum(occ3) + sum(occ4) + sum(occ5)
    ) as total
from agregateds
where subpath(geocode, 0, %d) = ?;
--@block test2
select edu,
    edu1,
    edu2,
    edu3,
    edu4,
    edu5,
    edu6,
    edu7,
    edu8,
    edu9,
    edu10,
    edu12,
    edu15,
    edu18,
    (
        edu1 + edu2 + edu3 + edu4 + edu5 + edu6 + edu7 + edu8 + edu9 + edu10 + edu12 + edu15 + edu18
    ) as total
from agregateds;
--@block
select edu,
    edu1,
    edu2
from agregateds;
--@block
select edu,
    edu1 as NoEducation,
    edu2 as Class1,
    edu3 as Class2,
    edu3 as Class3,
    edu4 as Class4,
    edu5 as Class5,
    edu6 as Class6,
    edu7 as Class7,
    edu8 as Class8,
    edu9 as Class9,
    edu10 as SCC,
    edu12 as HSC,
    edu15 as BachelorEquivalent,
    edu18 as MastersEquivalentOrHigher,
    (
        edu1 + edu2 + edu3 + edu4 + edu5 + edu6 + edu7 + edu8 + edu9 + edu10 + edu12 + edu15 + edu18
    ) as Total
from agregateds;
--@block
select sum(sex) as male,
    sum(sex2) as female,
    sum(sex3) as hijra,
    (sum(sex) + sum(sex2) + sum(sex3)) as total
from agregateds;
--@block
SELECT sum(hh_f) as Number_Of_Fishery_Household,
    (sum(hh_f) / sum(hh_sno)) * 100 as Percentage
from agregateds
where hh_f = 1;