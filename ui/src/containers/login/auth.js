import React from "react";
import { Redirect } from "react-router-dom";
import * as User from "../../utils/user";
import PropTypes from "prop-types";

const canAccess = roles => {
	for (let i = 0; i < roles.length; i++) {
		if (User.IsInRole(roles[i])) {
			return true;
		}
	}

	return false;
};

const propTypes = {
	roles: PropTypes.array,
	contentOnly: PropTypes.bool,
};
const defaultProps = {
	roles: [],
	contentOnly: false,
};

const Authorise = ({ children, roles, contentOnly }) => {
	const loginAction = contentOnly ? null : <Redirect to="/login" />;

	if (!User.IsAuthenticated()) {
		return loginAction;
	}

	if (roles.length > 0 && !canAccess(roles)) {
		return loginAction;
	}

	return children;
};

Authorise.propTypes = propTypes;
Authorise.defaultProps = defaultProps;

export default Authorise;
