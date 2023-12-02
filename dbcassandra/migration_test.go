package dbcassandra

import (
	"fmt"
	"github.com/biangacila/luvungula-go/global"
	"testing"
)

func TestMigration_RunFuel(t *testing.T) {
	var hub Migration
	hub.OriginHost = "pmis.easipath.com"
	hub.OriginDatabase = "fleetminder2"

	hub.DestinationHost = "voip.easipath.com"
	hub.DestinationDatabase = "fleetminder2"

	hub.ReplicationFactor = 1

	hub.Run()
	fmt.Println("qryDbCreate> ", hub.qryDbCreate)
	global.DisplayObject("tableCreateQuery", hub.tableCreateQuery)
	fmt.Println("total queries > ", len(hub.queries))

}
func TestMigration_Run(t *testing.T) {
	var hub Migration
	hub.OriginHost = "pmis.easipath.com"
	hub.OriginDatabase = "pmis"

	hub.DestinationHost = "voip.easipath.com"
	hub.DestinationDatabase = "pmis"

	hub.ReplicationFactor = 1

	hub.Run()
	fmt.Println("qryDbCreate> ", hub.qryDbCreate)
	global.DisplayObject("tableCreateQuery", hub.tableCreateQuery)
	fmt.Println("total queries > ", len(hub.queries))

}

/*

select appname,customernumber,displayname,name from company;
delete from company where appname='pmis' and customernumber='C100004';
delete from company where appname='pmis' and customernumber='C100005';

 appname | customernumber | displayname | name
---------+----------------+-------------+------------------------------------
    pmis |        C100006 |             |            Mshengu Toilet Hire C.C
    pmis |        C100007 |             | Cemex Trading Enterprise Pty (Ltd)
    pmis |        C100008 |             |          QSK Engineering (Pty) Ltd
    pmis |        C100015 |             |      Sivuyile Maintenance Services


update company set name='Client 1' where appname='pmis' and customernumber='C100006';
update company set name='Client 2' where appname='pmis' and customernumber='C100007';
update company set name='Client 3' where appname='pmis' and customernumber='C100008';
update company set name='Client 4' where appname='pmis' and customernumber='C100015';

*/
