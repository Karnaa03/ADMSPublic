# Agri Tracking

## Presentation

### Booklet state diagram

```puml
digraph fsm {
            "notRegistered" -> "registered" [ label = "register" ];
            "inCuttingStation" -> "inPreScanning" [ label = "moveToPreScanning" ];
            "inScanningStation" -> "archived" [ label = "moveToArchiveStation" ];
            "inIceBox" -> "archived" [ label = "moveToArchiveStation" ];
            "archived" -> "archived" [ label = "moveBookletFromBox" ];
            "inBatch" -> "registered" [ label = "removeFromBatch" ];
            "inBatch" -> "inCuttingStation" [ label = "moveToCuttingStation" ];
            "registered" -> "inBatch" [ label = "addInBatch" ];
            "inCuttingStation" -> "inIceBox" [ label = "freeze" ];
            "inPreScanning" -> "inScanningStation" [ label = "moveToScanStation" ];
            "registered" -> "notRegistered" [ label = "deregister" ];
        
            "inCuttingStation";
            "inPreScanning";
            "inScanningStation";
            "archived";
            "inIceBox";
            "inBatch";
            "notRegistered";
            "registered";
        }
```

### Crate state diagram

````puml
digraph fsm {
            "registered" -> "registered" [ label = "removeFromBatch" ];
            "registered" -> "inCuttingStation" [ label = "moveToCuttingStation" ];
            "registered" -> "registered" [ label = "addInBatch" ];
            "inCuttingStation" -> "inCuttingStation" [ label = "moveToCuttingStation" ];
            "inCuttingStation" -> "inPreScanning" [ label = "moveToPreScanning" ];
            "inPreScanning" -> "inPreScanning" [ label = "moveToPreScanning" ];
            "inScanningStation" -> "inScanningStation" [ label = "moveToScanStation" ];
            "inScanningStation" -> "inArchiveStation" [ label = "moveToArchiveStation" ];
            "inArchiveStation" -> "inArchiveStation" [ label = "moveToArchiveStation" ];
            "inPreScanning" -> "inScanningStation" [ label = "moveToScanStation" ];
            "inArchiveStation" -> "registered" [ label = "archive" ];
        
            "registered";
            "inCuttingStation";
            "inPreScanning";
            "inScanningStation";
            "inArchiveStation";
        }
````

### Shelf state diagram

```puml
digraph fsm {
            "registered" -> "inCuttingStation" [ label = "moveToCuttingStation" ];
            "registered" -> "registered" [ label = "addInBatch" ];
            "registered" -> "registered" [ label = "removeFromBatch" ];
            "inCuttingStation" -> "inPreScanning" [ label = "moveToPreScanning" ];
            "inScanningStation" -> "inArchiveStation" [ label = "moveToArchiveStation" ];
            "inArchiveStation" -> "inArchiveStation" [ label = "moveToArchiveStation" ];
            "inArchiveStation" -> "registered" [ label = "archive" ];
            "inCuttingStation" -> "inCuttingStation" [ label = "moveToCuttingStation" ];
            "inPreScanning" -> "inPreScanning" [ label = "moveToPreScanning" ];
            "inPreScanning" -> "inScanningStation" [ label = "moveToScanStation" ];
            "inScanningStation" -> "inScanningStation" [ label = "moveToScanStation" ];
        
            "registered";
            "inCuttingStation";
            "inPreScanning";
            "inScanningStation";
            "inArchiveStation";
        }
```

### Archive box state diagram

```puml
digraph fsm {
            "registered" -> "inUse" [ label = "moveToArchiveStation" ];
            "inUse" -> "inUse" [ label = "moveToArchiveStation" ];
            "inUse" -> "archived" [ label = "archive" ];
            "archived" -> "checkedOut" [ label = "checkOut" ];
            "checkedOut" -> "archived" [ label = "chekIn" ];
            "archived" -> "archived" [ label = "moveBookletFromBox" ];
        
            "registered";
            "inUse";
            "archived";
            "checkedOut";
        }
```

### FSM pre-post event

NewFSM constructs a FSM from events and callbacks.

The events and transitions are specified as a slice of Event structs specified as Events. Each Event is mapped to one or more internal transitions from Event.Src to Event.Dst.

Callbacks are added as a map specified as Callbacks where the key is parsed as the callback event as follows, and called in the same order:

1. `before_<EVENT>` - called before event named <EVENT>

2. `before_event` - called before all events

3. `leave_<OLD_STATE>` - called before leaving <OLD_STATE>

4. `leave_state` - called before leaving all states

5. `enter_<NEW_STATE>` - called after entering <NEW_STATE>

6. `enter_state` - called after entering all states

7. `after_<EVENT>` - called after event named <EVENT>

8. `after_event` - called after all events

There are also two short form versions for the most commonly used callbacks. They are simply the name of the event or state:

1. <NEW_STATE> - called after entering <NEW_STATE>

2. <EVENT> - called after event named <EVENT>

If both a shorthand version and a full version is specified it is undefined which version of the callback will end up in the internal map. This is due to the psuedo random nature of Go maps. No
checking for multiple keys is currently performed.

## Setup

The setup can be done via environment variables :

* `DB_HOST` : default "postgres.agritracking.svc.cluster.local:5432"
* `DB_USER` : default "agritracking"
* `DB_PASSWORD` : default "li7keegh4aexiToo"
* `DB_DATABASE` : default "agridev"
* `LISTEN_ADDR` : default "0.0.0.0:4000"
* `BASE_URL` : default "http://localhost:4000/"
* `OPENID_URL` : default "https://auth.solutions.im/auth/realms/solutions"
* `OPENID_CLIENT_ID` : default "test"
* `OPENID_CLIENT_SECRET` : default "0c862cea-64ca-4b07-b50f-2dca81a7a0b2"

## GeoCode generation

From the CSV version of Bernard's Excel file apply the following script :

```shell script
awk -F';' '{print $1";"$2";"$3";"$4";"$5";"$6";"$7";"$8";"$9";"$10";"$11";"$12";"$13";"$14";"$15";"$16";"$17";"$3"."$4"."$5"."$6"."$7"."$8 }' GeoCodes.csv
```

```shell script 
awk -F',' '{print $3"."$4"."$5"."$6"."$7"."$8";"$2";"$3";"$4";"$5";"$6";"$7";"$8";"$9";"$10";"$11";"$12";"$13";"$14";"$15";"$16}' agri_geo_bernard.csv | head
```

### Bernard view

```sql 
create or replace view bernard_booklets as
select b.number,
       ltree2text(b.geo_code_id) as geo_code_id,
       b.size,
       b.status,
       b.crate_number,
       b.archive_box_number,
       b.registered_on,
       b.registered_by,
       b.added_in_batch_on,
       b.added_in_batch_by,
       b.cut_on,
       b.cut_by,
       b.prepared_on,
       b.prepared_by,
       b.scanned_on,
       b.scanned_by,
       b.archived_on,
       b.archived_by,
       a.warehouse_row_number,
       a.warehouse_shelf_number,
       a.warehouse_shelf_level_number
from booklets b left join archive_boxes a on b.archive_box_number = a.number;
```        

## E-Flow view

```sql 
-- E-Flow view
create or replace view eflow as
    select b.number,
           b.geo_code_id,
           g.district,
           g.name_district,
           g.upazilla,
           g.name_upazilla,
           g."union",
           g.name_union,
           g.mouza,
           g.name_mouza,
           g.ca,
           g.name_counting_area,
           g.rmo,
           g.name_rmo
from booklets b left join geo_codes g on b.geo_code_id = g.geocode_id;

create user eflow password 'iemoh0Ig7Thoh9ho';
grant connect on database agritracking to eflow;
grant usage on schema public to eflow;
grant select on eflow to eflow;
```