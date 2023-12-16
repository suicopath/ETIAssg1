Create database carpool;

use carpool;
Create table passenger (passengerID INT NOT NULL PRIMARY KEY, fisrtName Varchar(50), lastName varchar(30), phoneNo varchar(20), email varchar (50), carOwner Bool);

Create table carowner (passengerID INT, carPlateNo varchar(50) NOT NULL, licenseNo Varchar(50), maxPassenger int, carColor varchar(15), yearRel int);

Create table trip (tripID INT NOT NULL PRIMARY KEY, pickupAddr varchar(50), alterPickupAddr varchar(50), startTrip varchar(20), destinationAddr varchar(50), maxPassenger INT, passengerID INT);


INSERT INTO passenger VALUES (1, 'Goh', 'Chok Tong', '88237164', 'ctgoh@hotmail.com', FALSE),
(2, 'Lim', 'Yew Hock', '88732192', 'yhlim@hotmail.com', FALSE),
(3, 'Heng', 'Swee Keat', '90293124', 'skheng@hotmail.com', TRUE);

INSERT INTO carowner VALUES (3, 'SGX39162Z', '8791274905', 4, 'Green', 2018);

INSERT INTO trip VALUES (1 ,'Yishun St 72 Blk 749', 'Yishun Central Blk 226', CURRENT_TIMESTAMP, '52 International Rd Blk 3', 4, 3);
INSERT INTO trip VALUES (2 ,'Yishun Ring Rd Blk 987', 'Canberra Lane Ave 2', CURRENT_TIMESTAMP, 'Shenton Way Blk 43', 6, 2);


