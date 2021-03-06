package scrape

//SubjectsMap is a map of the start of the sds code and the subject
//Kent has no easy way to fetch this data, so I'm collecting it from a <select> tag in this page:
//https://www.kent.ac.uk/courses/modules
var SubjectsMap = map[string]string{
	"AC":   "Accounting & Finance",
	"SE":   "Anthropology",
	"LABS": "Apprenticeships: Laboratory Scientists",
	"AR":   "Architecture",
	"ART":  "Arts and Media",
	"BI":   "Biosciences",
	"CB":   "Business",
	"CH":   "Chemistry",
	"CL":   "Classical & Archaeological Studies",
	"CP":   "Comparative Literary Studies",
	"CO":   "Computing, Computer Science",
	"DI":   "Conservation",
	"DR":   "Drama",
	"EC":   "Economics",
	"EL":   "Engineering and Digital Arts",
	"LL":   "English Language and Linguistics",
	"EN":   "English, American and Postcolonial Literatures",
	"CR":   "Event and Experience Design",
	"FI":   "Film Studies",
	"FA":   "Fine Art",
	"PS":   "Forensic Science",
	"FR":   "French",
	"GE":   "German",
	"HM":   "Heritage Management",
	"LS":   "Hispanic Studies",
	"HI":   "History",
	"HA":   "History & Philosophy of Art",
	"HU":   "Humanities",
	"LZ":   "International Foundation Programme",
	"IT":   "Italian",
	"JN":   "Journalism",
	"LW":   "Law",
	"MA":   "Mathematics, Statistics and Actuarial Science ",
	"MT":   "Medieval Studies",
	"MU":   "Music",
	"CMAT": "Music Technology",
	"PHAM": "Pharmacy",
	"PL":   "Philosophy",
	"PH":   "Physics",
	"PO":   "Politics and International Relations",
	"WL":   "Professional Practise",
	"SP":   "Psychology",
	"TH":   "Religious Studies",
	"SA":   "Social Policy",
	"SO":   "Sociology",
	"SS":   "Sport",
	"TZ":   "Tizard Centre",
	"UN":   "UELT",
}
