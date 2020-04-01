import React from "react";
import { Header } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import ListContainer from "../../containers/plugins/list";

const List = () => (
	<Authorise roles={["Admin", "Plugin Manager"]}>
		<Header as="h1">Plugins</Header>

		<p>View and manage all plugins.</p>

		<ListContainer />
	</Authorise>
);

export default List;
