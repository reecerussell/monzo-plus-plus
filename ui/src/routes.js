import React from "react";

const UsersList = React.lazy(() => import("./views/users/list"));
const UsersDetails = React.lazy(() => import("./views/users/details"));
const UsersPending = React.lazy(() => import("./views/users/pending"));

const RolesList = React.lazy(() => import("./views/roles/list"));
const RolesCreate = React.lazy(() => import("./views/roles/create"));
const RolesEdit = React.lazy(() => import("./views/roles/edit"));

const PluginsList = React.lazy(() => import("./views/plugins/list"));

const AccountIndex = React.lazy(() => import("./views/account/account"));
const AccountChangePassword = React.lazy(() =>
	import("./views/account/changePassword")
);
const AccountPlugins = React.lazy(() => import("./views/account/plugins"));
const AccountDelete = React.lazy(() => import("./views/account/delete"));
const AccountSetAccount = React.lazy(() =>
	import("./views/account/setAccount")
);

const Login = React.lazy(() => import("./views/login/login"));
const Register = React.lazy(() => import("./views/login/register"));

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
		name: "roles",
		path: "/roles",
		exact: true,
		component: RolesList,
	},
	{
		name: "create role",
		path: "/roles/create",
		component: RolesCreate,
	},
	{
		name: "edit role",
		path: "/roles/edit/:id",
		component: RolesEdit,
	},
	{
		name: "plugins",
		path: "/plugins",
		component: PluginsList,
		exact: true,
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
		component: AccountChangePassword,
	},
	{
		name: "account plugins",
		path: "/account/plugins",
		component: AccountPlugins,
	},
	{
		name: "delete account",
		path: "/account/delete",
		component: AccountDelete,
	},
	{
		name: "set account",
		path: "/setAccount",
		exact: true,
		component: AccountSetAccount,
	},
	{
		name: "login",
		path: "/login",
		component: Login,
		exact: true,
	},
	{
		name: "register",
		path: "/register",
		component: Register,
		exact: true,
	},
];
