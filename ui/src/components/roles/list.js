import React from "react";
import {
	Table,
	Message,
	Loader,
	Segment,
	Form,
	Button,
} from "semantic-ui-react";
import { Link } from "react-router-dom";
import PropTypes from "prop-types";
import * as User from "../../utils/user";

const propTypes = {
	roles: PropTypes.array,
	loading: PropTypes.bool.isRequired,
	error: PropTypes.string,
	searchTerm: PropTypes.string,
	updateSearchTerm: PropTypes.func,
	onSearch: PropTypes.func,
	handleFlush: PropTypes.func.isRequired,
};
const defaultProps = {
	roles: [],
	error: null,
};

const List = ({
	roles,
	loading,
	error,
	searchTerm,
	updateSearchTerm,
	onSearch,
	handleFlush,
}) => {
	const errorMessage = error ? (
		<Message error header="An error occured!" content={error} />
	) : null;

	const getName = name => {
		if (searchTerm.length < 1) {
			return name;
		}

		const regex = new RegExp(name, "gi");

		return name.replace(
			regex,
			str => "<span style='background-color: yellow;'>" + str + "</span>"
		);
	};

	const rows =
		roles && roles.length > 0 ? (
			roles.map((role, idx) => (
				<Table.Row key={idx}>
					<Table.Cell>
						<span
							dangerouslySetInnerHTML={{
								__html: getName(role.name),
							}}
						/>
					</Table.Cell>
					<Table.Cell>
						<Link to={"/roles/edit/" + role.id}>View</Link>
					</Table.Cell>
				</Table.Row>
			))
		) : (
			<Table.Row>
				<Table.Cell colSpan="3">
					<b>
						{searchTerm.length > 0
							? "No roles were found for this search, try changing your query."
							: "No roles exist."}
					</b>
				</Table.Cell>
			</Table.Row>
		);

	return (
		<>
			{errorMessage}
			<Segment basic>
				<Loader active={loading} />

				<Form onSubmit={onSearch}>
					<Form.Group>
						<Form.Input
							placeholder="Search..."
							value={searchTerm}
							onChange={updateSearchTerm}
						/>
						<Form.Button content="Search" />
						{User.IsInRole("Role Manager") ? (
							<Form.Field>
								<Button
									color="link"
									type="button"
									onClick={handleFlush}
								>
									Flush Permissions
								</Button>
							</Form.Field>
						) : null}
					</Form.Group>
				</Form>

				<Table striped>
					<Table.Header>
						<Table.Row>
							<Table.HeaderCell>Name</Table.HeaderCell>
							<Table.HeaderCell></Table.HeaderCell>
						</Table.Row>
					</Table.Header>

					<Table.Body>{rows}</Table.Body>
				</Table>
			</Segment>
		</>
	);
};

List.propTypes = propTypes;
List.defaultProps = defaultProps;

export default List;
