
ALTER TABLE documents ADD hash   TEXT;
ALTER TABLE documents ADD signature  TEXT;

update documents
	set hash = 'A6536dB0989867083a47EA5344cEa382b0Bf4F21'
	  , signature = 'eA4afF1497E5763A9980589e1b1C40018A9BCb7dC8511cB5f0632330906844D58A7F611F777CC4A80'
;