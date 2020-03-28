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
					},
					"403": {
						description: "Insufficient permissions",
					},
					"404": {
						description: "No user with the id was found.",
					},
					"200": {
						description: "Successful request.",
					},
				},
			},
		},
	},
	securityDefinitions: {
		petstore_auth: {
			type: "oauth2",
			authorizationUrl: "http://petstore.swagger.io/oauth/dialog",
			flow: "implicit",
			scopes: {
				"write:pets": "modify pets in your account",
				"read:pets": "read your pets",
			},
		},
		api_key: {
			type: "apiKey",
			name: "api_key",
			in: "header",
		},
	},
	definitions: {
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
	},
	externalDocs: {
		description: "View source",
		url: "http://github.com/reecerussell/monzo-plus-plus",
	},
};
