package main

var routes = Routes{
	Route{
		"BuildList",
		"GET",
		"/api/buildings",
		BuildList,
	},
	Route{
		"BuildById",
		"GET",
		"/api/buildings/{id}",
		BuildById,
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
