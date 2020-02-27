import React from "react";
import { Card, Icon, Message } from "semantic-ui-react";
import PropTypes from "prop-types";

const propTypes = {
	username: PropTypes.string,
	enabled: PropTypes.bool,
	roles: PropTypes.array,
	loading: PropTypes.bool,
	error: PropTypes.string,
};
const defaultProps = {
	roles: [],
};

const Details = ({ username, enabled, roles, loading, error }) => {
	if (loading) {
		return (
			<Card>
				<Card.Content header={username || "User"} />
				<Card.Content description="Loading..." />
				<Card.Content extra>
					<Icon name="user" />
					Loading...
				</Card.Content>
			</Card>
		);
	}

	if (error) {
		return <Message error header="An error occured!" content={error} />;
	}

	return (
		<Card>
			<Card.Content header={username} />
			<Card.Content
				description={`${username} is ${enabled ? "" : "not "} enabled.`}
			/>
			<Card.Content extra>
				<Icon name="user" />
				{roles.length} Roles
			</Card.Content>
		</Card>
	);
};

Details.propTypes = propTypes;
Details.defaultProps = defaultProps;

export default Details;
