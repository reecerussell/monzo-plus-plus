import React from "react";
import { Header, Grid } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import Layout from "../../components/account/layout";
import PluginsContainer from "../../containers/account/plugins";

const Plugins = () => {
	return (
		<Authorise>
			<Layout>
				<Header as="h2">Plugins</Header>

				<Grid stackable>
					<Grid.Column width="12">
						<p>
							Plugins are individual features that you can add to
							your account. You can add as many as you like,
							simpily by enabling one below.
						</p>
						<p>
							If you ever decide you'd like to stop one, either
							temporarily or permanently, you can click the
							"Disable Plugin" button in the section below.
						</p>
					</Grid.Column>
				</Grid>

				<PluginsContainer />
			</Layout>
		</Authorise>
	);
};

export default Plugins;
