import React from "react";
import { Header } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import ListContainer from "../../containers/roles/list";
import { Link } from "react-router-dom";

const List = () => (
	<Authorise roles={["Admin", "Role Manager"]}>
		<Header as="h1">Roles</Header>

		<p>
			<Link to="/roles/create">Create a new role</Link>
		</p>

		<ListContainer />
	</Authorise>
);

export default List;
