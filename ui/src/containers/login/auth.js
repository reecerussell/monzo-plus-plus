import React from "react";
import { Redirect, useLocation } from "react-router-dom";
import * as User from "../../utils/user";
import PropTypes from "prop-types";

const canAccess = (roles) => {
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
	const { pathname } = useLocation();
	let canProceed = true;

	if (!User.IsAuthenticated()) {
		canProceed = false;
	}

	if (roles.length > 0 && !canAccess(roles)) {
		canProceed = false;
	}

	if (!canProceed) {
		if (contentOnly) {
			return null;
		}

		return <Redirect to="/login" />;
	}

	return children;
};

Authorise.propTypes = propTypes;
Authorise.defaultProps = defaultProps;

export default Authorise;
