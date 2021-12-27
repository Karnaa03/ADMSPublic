package model

import (
	"time"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"
)

// Events
const register = "register"
const deregister = "deregister"
const addInBatch = "addInBatch"
const removeFromBatch = "removeFromBatch"
const moveToCuttingStation = "moveToCuttingStation"
const moveToPreScanning = "moveToPreScanning"
const moveToScanStation = "moveToScanStation"
const moveToArchiveStation = "moveToArchiveStation"
const moveBookletFromBox = "moveBookletFromBox"
const archive = "archive"
const checkOut = "checkOut"
const checkIn = "chekIn"
const freeze = "freeze"

// States
const notRegistered = "notRegistered"
const registered = "registered"
const inBatch = "inBatch"
const inCuttingStation = "inCuttingStation"
const inPreScanning = "inPreScanning"
const inScanningStation = "inScanningStation"
const inArchiveStation = "inArchiveStation"
const checkedOut = "checkedOut"
const archived = "archived"
const iceBox = "inIceBox"
const inUse = "inUse"

type Event struct {
	Id               int64 `pg:",pk" binding:"-"`
	SourceState      string
	EventType        string
	FinalState       string
	TimeStamp        time.Time
	IdentityId       string
	Identity         *gin_oidc.Identity `pg:"rel:has-one"`
	Subject          EventSubject       `pg:"-"`
	BookletNumber    string
	Booklet          *Booklet `pg:"rel:has-one"`
	CrateNumber      string
	Crate            *Crate `pg:"rel:has-one"`
	ShelfNumber      string
	Shelf            *Shelf `pg:"rel:has-one"`
	ArchiveBoxNumber string
	ArchiveBox       *ArchiveBox `pg:"rel:has-one"`
}

type EventSubject interface {
	GetType() string
	GetId() string
	// GetActor() string
}
