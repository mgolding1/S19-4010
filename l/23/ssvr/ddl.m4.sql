
m4_include(setup.m4)

CREATE SEQUENCE t_doc_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;

CREATE TABLE "documents" (
	  "id"					m4_uuid_type() DEFAULT uuid_generate_v4() not null primary key
	, "user_id"				m4_uuid_type()
	, "document_hash" 		TEXT	-- is this the same as "hash" below?
	, "email" 				TEXT
	, "real_name" 			TEXT				-- should have a soundex type key on this. and a word-index key
	, "phone_number" 		TEXT
	, "address_usps" 		TEXT				-- shlold have a word-index type key
	, "document_file_name"	TEXT
	, "file_name"			TEXT
	, "orig_file_name"		TEXT
	, "url_file_name"		TEXT
	, "hash"				char varying(32)
	, "signature"			char varying(68)
	, "txid"				TEXT				-- transaction hash that is repoted on request to Eth
	, "ethstatus"			text				-- status of eth transaction
	, "blockhash"			TEXT				-- block hash where posted
	, "blockno"				TEXT				-- block number where posted
	, "blockerr"			TEXT				-- error if tranaction failed
	, "note" 				TEXT
	, "qr_id" 				TEXT
	, "qr_enc_id"			TEXT
	, "qr_data" 			TEXT
	, "qr_to_url" 			TEXT
	, "updated" 			timestamp
	, "created" 			timestamp default current_timestamp not null
);

create index "documents_p1" on "documents" ( "hash" );
create index "documents_p2" on "documents" ( "ethstatus" );
create index "documents_p3" on "documents" ( "email" );
create index "documents_p4" on "documents" ( "real_name" );
create index "documents_p5" on "documents" ( "phone_number" );
create index "documents_p6" on "documents" ( "orig_file_name" );

m4_updTrig(documents)

CREATE TABLE "users" (
	  "id"					m4_uuid_type() DEFAULT uuid_generate_v4() not null primary key
	, "user_id"				m4_uuid_type()
	, "username" 			TEXT
	, "password"			TEXT
	, "updated" 			timestamp
	, "created" 			timestamp default current_timestamp not null
);

m4_updTrig(users)


CREATE TABLE "states" (
	  "id"					m4_uuid_type() DEFAULT uuid_generate_v4() not null primary key
	, "state"	 			TEXT
	, "final_st"			TEXT
	, "final"				TEXT
	, "action_call"			TEXT
	, "desc"				TEXT
	, "updated" 			timestamp
	, "created" 			timestamp default current_timestamp not null
);

m4_updTrig(states)

--		x. Add "ItemEvent" - for tracking additional documents/events related to.... (Demo for Hash Folks)
--			Event 		Ev-Final  	Final	Desc
--			--------	------		-----	---------------------------
--			Create		yes			no		Enable QR code to be live for apartiular user. -- Assoc. QR with a "originator".
--			QualityTest	yes			no		Add a mesurment/quality document - final test result.
--			AddDocument no			no		Associated tracking document.
--			AddText		no			no		Associated tracking text.
--			Processing	no			no		Process into new product.
--			ShipingBeg	no			no		Start of Shipping.
--			ShipingEnd	no			no		End of Shipping.
--			Location    no			no		Location / date time stamp.
--			Split		no			no		Split into N new tags.   Tag array returned.
--			Combine		no			no		Combine a set of M tags into a new tag. (Mix products together)
--			Delete		yes			yes		Item is deleted - completely final.
--			Inspection	no			no		Inspection occured - Inspecotr ID info etc.
--			FarmInfo	no			no		Add info about where crop is raised.
--			DataExport	yes			yes		Data moved to different network (IFT for example).
--			CustAcc		no			no		Customer Access of Data (QR Code Access).
--			EndUserSale yes			yes		Sold to an end consumer.

