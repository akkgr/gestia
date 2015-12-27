package main

var routes = Routes{
	Route{
		"BuildIndex",
		"GET",
		"/api/buildings",
		BuildIndex,
	},
	Route{
		"BuildShow",
		"GET",
		"/api/buildings/{id}",
		BuildShow,
	},
	Route{
		"BuildInsert",
		"POST",
		"/api/buildings",
		BuildInsert,
	},
	Route{
		"BuildUpdate",
		"PUT",
		"/api/buildings/{id}",
		BuildUpdate,
	},
	Route{
		"BuildDelete",
		"DELETE",
		"/api/buildings/{id}",
		BuildDelete,
	},
}
