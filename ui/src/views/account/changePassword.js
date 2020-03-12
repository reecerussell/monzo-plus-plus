import React from "react";
import { Authorise } from "../../containers/login";
import ChangePasswordContainer from "../../containers/account/changePassword";

const ChangePassword = () => (
	<Authorise>
		<ChangePasswordContainer />
	</Authorise>
);

export default ChangePassword;
