import React from "react";
import PropTypes from "prop-types";
import { Button, Modal, Header, Icon } from "semantic-ui-react";
import { Redirect } from "react-router-dom";

const propTypes = {
	handleSubmit: PropTypes.func.isRequired,
	isModalOpen: PropTypes.bool,
	toggleModal: PropTypes.func.isRequired,
	loading: PropTypes.bool,
	redirect: PropTypes.string,
};
const defaultProps = {
	isModalOpen: false,
	loading: false,
	redirect: null,
};

const Delete = ({
	handleSubmit,
	isModalOpen,
	toggleModal,
	loading,
	redirect,
}) => {
	if (redirect) {
		return <Redirect to={redirect} />;
	}

	const button = (
		<Button color="red" loading={loading} onClick={toggleModal}>
			Delete
		</Button>
	);

	return (
		<Modal
			trigger={button}
			open={isModalOpen}
			onClose={toggleModal}
			size="small"
			basic
		>
			<Header icon="trash" content="Delete" />
			<Modal.Content>
				<h3>Are you sure you'd like to delete this role?</h3>
			</Modal.Content>
			<Modal.Actions>
				<Button onClick={toggleModal}>
					<Icon name="close" /> Cancel
				</Button>
				<Button color="red" onClick={handleSubmit}>
					<Icon name="checkmark" /> Delete
				</Button>
			</Modal.Actions>
		</Modal>
	);
};

Delete.propTypes = propTypes;
Delete.defaultProps = defaultProps;

export default Delete;
