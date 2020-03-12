import React from "react";
import { Authorise } from "../../containers/login";
import IndexContainer from "../../containers/account/index";

const Account = () => (
	<Authorise>
		<IndexContainer />
	</Authorise>
);

export default Account;
