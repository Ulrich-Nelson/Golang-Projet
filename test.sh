go run mdb.go clear
go run mdb.go init

go run mdb.go create-practitioner -practitionerID=ff270805549658d057bb4ea15dc77303 -username=vidal -password=pipo

go run mdb.go create-patient -patientID=dca5b0ded732c3a3a94c34f45ceeada1 -practitionerID=ff270805549658d057bb4ea15dc77303 -data='{"firstname":"Jean","lastname":"Amand","birthdate":"1960"}' -password=pipo
go run mdb.go create-patient -patientID=c645cd076995ea9a6d28cb0c9bbd8e14 -practitionerID=ff270805549658d057bb4ea15dc77303 -data='{"firstname":"Sandrine","lastname":"Delage","birthdate":"1953"}' -password=pipo
go run mdb.go create-patient -patientID=577d9a886540e20915d7a6494d5f761f -practitionerID=ff270805549658d057bb4ea15dc77303 -data='{"firstname":"Alphonse","lastname":"Ferreira","birthdate":"1929"}' -password=pipo

go run mdb.go create-protocol -protocolName=WalkData -real=T -real=Ox -real=Oy -real=Oz -real=Ax -real=Ay -real=Az
go run mdb.go create-experiment -experimentID=8B8F3111675CACEF32B188D76AB85294 -packetID=c29f0d94e1078c1b64040f1db3a2ea5a -practitionerID=ff270805549658d057bb4ea15dc77303 -patientID=dca5b0ded732c3a3a94c34f45ceeada1
go run mdb.go create-experiment -experimentID=0C37B5B6AB705643D6E6C6CC63CF5A7D -packetID=c29f0d94e1078c1b64040f1db3a2ea5a -practitionerID=ff270805549658d057bb4ea15dc77303 -patientID=c645cd076995ea9a6d28cb0c9bbd8e14
go run mdb.go create-experiment -experimentID=C2DF8E9236BDC58112A73720D8EB9DEF -packetID=c29f0d94e1078c1b64040f1db3a2ea5a -practitionerID=ff270805549658d057bb4ea15dc77303 -patientID=c645cd076995ea9a6d28cb0c9bbd8e14
go run mdb.go create-experiment -experimentID=8A09507313230A48B78717DD439A63F1 -packetID=c29f0d94e1078c1b64040f1db3a2ea5a -practitionerID=ff270805549658d057bb4ea15dc77303 -patientID=577d9a886540e20915d7a6494d5f761f

go run mdb.go inject -protocolName=WalkData -filename=data/WL_Amand_2019_04_23_12_18_17.csv
go run mdb.go inject -protocolName=WalkData -filename=data/WL_Delage_2019_04_23_12_58_51.csv
go run mdb.go inject -protocolName=WalkData -filename=data/WL_Delage_Sandrine_2019_04_23_01_20_43.csv
go run mdb.go inject -protocolName=WalkData -filename=data/WL_Ferreira_2019_04_23_10_25_10.csv

go run mdb.go create-protocol -protocolName=WalkResult -text=direction -real=duration -real=time_start -real=yaw_start -real=time_stop -real=yaw_stop -real=velocity_ratio -real=velocity_averaging
go run mdb.go inject -protocolName=WalkResult -filename=data/summary_db.csv

go run mdb.go create-patient -patientID=8fda0c3ae1310dafb431d0c8d9d3549f -practitionerID=ff270805549658d057bb4ea15dc77303 -data='{"firstname":"Mohammed","lastname":"Abrami","birthdate":"1961"}' -password=pipo
go run mdb.go create-patient -patientID=334785e432c56f0cd6b18f04a23a1888 -practitionerID=ff270805549658d057bb4ea15dc77303 -data='{"firstname":"Jules","lastname":"Barthelemy","birthdate":"1952"}' -password=pipo

go run mdb.go create-protocol -protocolName=WalkData2 -real=T -real=Ox -real=Oy -real=Oz -real=Ax -real=Ay -real=Az
go run mdb.go create-experiment -experimentID=C923E072FB870A0181DDB2A8FAC76D12 -practitionerID=ff270805549658d057bb4ea15dc77303 -patientID=8fda0c3ae1310dafb431d0c8d9d3549f
go run mdb.go create-experiment -experimentID=A8BA22F208EA357F2DA6367193694494 -practitionerID=ff270805549658d057bb4ea15dc77303 -patientID=334785e432c56f0cd6b18f04a23a1888

go run mdb.go inject -protocolName=WalkData2 -filename=data/WR_Abrami_2019_06_11_12_29_51.csv
go run mdb.go inject -protocolName=WalkData2 -filename=data/WR_Barthelemy_2019_06_11_10_33_45.csv

go run mdb.go create-protocol -protocolName=WalkResult2 -text=direction -real=duration -real=time_start -real=yaw_start -real=time_stop -real=yaw_stop -real=velocity_ratio -real=velocity_averaging
go run mdb.go inject -protocolName=WalkResult2 -filename=data/summary_db2.csv







