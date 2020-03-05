import React from "react";

const UsersList = React.lazy(() => import("./views/users/list"));
const UsersDetails = React.lazy(() => import("./views/users/details"));
const UsersPending = React.lazy(() => import("./views/users/pending"));

const AccountIndex = React.lazy(() => import("./views/account/account"));

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
	{
		name: "account",
		path: "/account",
		component: AccountIndex,
		exact: true,
	},
	{
		name: "change password",
		path: "/account/changepassword",
		component: AccountIndex,
	},
	{
		name: "account plugins",
		path: "/account/plugins",
		component: AccountIndex,
	},
	{
		name: "delete account",
		path: "/account/delete",
		component: AccountIndex,
	},
];
