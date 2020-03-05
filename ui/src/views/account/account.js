import React from "react";
import { Header } from "semantic-ui-react";
import IndexContainer from "../../containers/account/index";

const Account = () => (
	<>
		<Header as="h1">Account</Header>
		<p>Information about your account.</p>
		<IndexContainer />
	</>
);

export default Account;
