import React from "react";
import { Header } from "semantic-ui-react";
import { Authorise } from "../../containers/login";
import Layout from "../../components/account/layout";
import PluginsContainer from "../../containers/account/plugins";

const Plugins = () => {
	return (
		<Authorise>
			<Layout>
				<Header as="h2">Plugins</Header>

				<p>
					Plugins are individual features that you can add to your
					account. You can add as many as you like, simpily by
					enabling one below.
				</p>

				<PluginsContainer />
			</Layout>
		</Authorise>
	);
};

export default Plugins;
