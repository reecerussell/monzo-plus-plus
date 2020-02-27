import React from "react";

const UsersList = React.lazy(() => import("./views/users/list"));
const UsersDetails = React.lazy(() => import("./views/users/details"));
const UsersPending = React.lazy(() => import("./views/users/pending"));

export default [
	{
		name: "users",
		path: "/users",
		exact: true,
		component: UsersList,
	},
	{
		name: "pending users",
		path: "/users/pending",
		component: UsersPending,
	},
	{
		name: "user details",
		path: "/users/details/:id",
		component: UsersDetails,
	},
];
