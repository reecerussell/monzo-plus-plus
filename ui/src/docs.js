export default {
	swagger: "2.0",
	info: {
		description: "Monzo++",
		version: "1.0.0",
		title: "Monzo++",
		license: {
			name: "MIT",
			url:
				"https://github.com/reecerussell/monzo-plus-plus/blob/development/LICENSE",
		},
	},
	host: "mpp.reece-russell.co.uk",
	basePath: "/api",
	tags: [
		{
			name: "Token",
			description: "Endpoints for the OAuth flow.",
		},
		{
			name: "Users",
			description: "Endpoints to manage users.",
		},
		// {
		// 	name: "user",
		// 	description: "Operations about user",
		// 	externalDocs: {
		// 		description: "Find out more about our store",
		// 		url: "http://swagger.io",
		// 	},
		// },
	],
	schemes: ["https"],
	paths: {
		"/auth/token": {
			post: {
				tags: ["Token"],
				summary: "Endpoints for the OAuth flow.",
				description:
					"Returns an OAuth access token, used to authorise and authenticate API requests, in exchange for a username and password.",
				operationId: "authToken",
				consumes: ["application/json"],
				produces: ["application/json"],
				parameters: [
					{
						in: "body",
						name: "body",
						description: "User's username and password.",
						required: true,
						schema: {
							$ref: "#/definitions/UserCredential",
						},
					},
				],
				responses: {
					"400": {
						description: "Invalid input data.",
					},
					"200": {
						description: "Successful request.",
					},
				},
			},
		},
		"/auth/users/{id}": {
			get: {
				tags: ["Users"],
				summary: "Returns a specific user.",
				description:
					"An individual user record with the given id is returned.",
				operationId: "singleUser",
				produces: ["application/json"],
				parameters: [
					{
						in: "id",
						name: "id",
						description: "A user's id.",
						required: true,
					},
				],
				responses: {
					"401": {
						description: "Unauthorized request.",
						schema: {
							$ref: "#/definitions/Error",
						},
					},
					"403": {
						description: "Insufficient permissions",
						schema: {
							$ref: "#/definitions/Error",
						},
					},
					"404": {
						description: "No user with the id was found.",
						schema: {
							$ref: "#/definitions/Error",
						},
					},
					"200": {
						description: "Successful request.",
						schema: {
							$ref: "#/definitions/User",
						},
					},
				},
			},
			delete: {
				tags: ["Users"],
				summary: "Deletes a user.",
				description: "Deletes an individual user with the given id.",
				operationId: "deleteUser",
				produces: ["application/json"],
				parameters: [
					{
						in: "id",
						name: "id",
						description: "A user's id.",
						required: true,
					},
				],
				responses: {
					"401": {
						description: "Unauthorized request.",
						schema: {
							$ref: "#/definitions/Error",
						},
					},
					"403": {
						description: "Insufficient permissions",
						schema: {
							$ref: "#/definitions/Error",
						},
					},
					"404": {
						description: "No user with the id was found.",
						schema: {
							$ref: "#/definitions/Error",
						},
					},
					"200": {
						description: "Successful request.",
					},
				},
			},
		},
		"/auth/users": {
			get: {
				tags: ["Users"],
				summary: "A list of users.",
				description:
					"Returns a list of users, matching the search term (if any).",
				operationId: "listUser",
				produces: ["application/json"],
				parameters: [
					{
						in: "query",
						name: "term",
						description: "A search term",
						required: false,
					},
				],
				responses: {
					"401": {
						description: "Unauthorized request.",
						schema: {
							$ref: "#/definitions/Error",
						},
					},
					"403": {
						description: "Insufficient permissions",
						schema: {
							$ref: "#/definitions/Error",
						},
					},
					"200": {
						description: "Successful request.",
						schema: {
							type: "array",
							items: {
								$ref: "#/definitions/User",
							},
						},
					},
				},
			},
			post: {
				tags: ["Users"],
				summary: "Create a user",
				description:
					"Uses the data in the request body to create a new user.",
				operationId: "authToken",
				consumes: ["application/json"],
				produces: ["application/json"],
				parameters: [
					{
						in: "body",
						name: "body",
						description: "New user data.",
						required: true,
						schema: {
							$ref: "#/definitions/UserCredential",
						},
					},
				],
				responses: {
					"400": {
						description: "Invalid input data.",
					},
					"200": {
						description: "Successful request.",
					},
				},
			},
		},
	},
	securitySchemes: {
		BearerAuth: {
			type: "http",
			scheme: "bearer",
			bearerFormat: "JWT",
		},
	},
	definitions: {
		Error: {
			type: "object",
			properties: {
				error: {
					type: "string",
				},
			},
		},
		UserCredential: {
			type: "object",
			properties: {
				username: {
					type: "string",
				},
				password: {
					type: "string",
				},
			},
		},
		UpdateUser: {
			type: "object",
			properties: {
				id: {
					type: "string",
				},
				username: {
					type: "string",
				},
			},
		},
		User: {
			type: "object",
			properties: {
				id: {
					type: "string",
				},
				username: {
					type: "string",
				},
				enabled: {
					type: "boolean",
				},
				dateEnabled: {
					type: "string",
					nullable: true,
				},
			},
		},
		UserPlugin: {
			type: "object",
			properties: {
				userId: {
					type: "string",
				},
				pluginId: {
					type: "string",
				},
			},
		},
		UserRole: {
			type: "object",
			properties: {
				userId: {
					type: "string",
				},
				roleId: {
					type: "string",
				},
			},
		},
	},
	externalDocs: {
		description: "View source",
		url: "http://github.com/reecerussell/monzo-plus-plus",
	},
};
