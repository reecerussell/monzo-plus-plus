import React from "react";
import { Header } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import ListContainer from "../../containers/users/list";
import { Link } from "react-router-dom";

const List = () => (
	<Authorise roles={["Admin", "User Manager"]}>
		<Header as="h1">Users</Header>

		<p>
			<Link to="/users/pending">View pending users</Link>
		</p>

		<ListContainer />
	</Authorise>
);

export default List;
