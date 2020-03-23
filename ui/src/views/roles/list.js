import React from "react";
import { Header } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import ListContainer from "../../containers/roles/list";

const List = () => (
	<Authorise roles={["Admin"]}>
		<Header as="h1">Roles</Header>

		<ListContainer />
	</Authorise>
);

export default List;
