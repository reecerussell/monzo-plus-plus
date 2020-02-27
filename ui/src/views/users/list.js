import React from "react";
import { Header } from "semantic-ui-react";
import ListContainer from "../../containers/users/list";
import { Link } from "react-router-dom";

const List = () => (
	<>
		<Header as="h1">Users</Header>

		<p>
			<Link to="/users/pending">View pending users</Link>
		</p>

		<ListContainer />
	</>
);

export default List;
