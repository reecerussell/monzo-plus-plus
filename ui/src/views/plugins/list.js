import React from "react";
import { Header, Grid } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import ListContainer from "../../containers/plugins/list";

const List = () => (
	<Authorise roles={["Admin", "Plugin Manager"]}>
		<Header as="h1">Plugins</Header>

		<p>View and manage all plugins.</p>

		<Grid stackable>
			<Grid.Column width="8">
				<ListContainer />
			</Grid.Column>
		</Grid>
	</Authorise>
);

export default List;
