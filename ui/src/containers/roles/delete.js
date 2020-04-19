import React, { useState } from "react";
import PropTypes from "prop-types";
import { Fetch } from "../../utils/fetch";
import Delete from "../../components/roles/delete";

const propTypes = {
	id: PropTypes.string.isRequired,
	onError: PropTypes.func.isRequired,
};
const defaultProps = {};

const DeleteContainer = ({ id, onError }) => {
	const [loading, setLoading] = useState(false);
	const [redirect, setRedirect] = useState(null);
	const [isModelOpen, setIsModalOpen] = useState(false);

	const toggleModal = (e) => {
		e.preventDefault();
		setIsModalOpen(!isModelOpen);
	};

	const handleSubmit = async (e) => {
		e.preventDefault();

		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			`api/auth/roles/${id}`,
			{
				method: "DELETE",
			},
			async () => setRedirect("/roles"),
			(err) => {
				onError(err);
				setIsModalOpen(false);
			}
		);

		setLoading(false);
	};

	return (
		<Delete
			loading={loading}
			toggleModal={toggleModal}
			handleSubmit={handleSubmit}
			isModalOpen={isModelOpen}
			redirect={redirect}
		/>
	);
};

DeleteContainer.propTypes = propTypes;
DeleteContainer.defaultProps = defaultProps;

export default DeleteContainer;
