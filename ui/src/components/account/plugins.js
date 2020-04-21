import React, { Suspense } from "react";
import {
	Item,
	Button,
	Icon,
	Form,
	Message,
	Divider,
	Label,
	ButtonGroup,
} from "semantic-ui-react";
import * as User from "../../utils/user";

const Preferences = React.lazy(() =>
	import("../../containers/plugins/budget/preferences")
);

const getPluginModal = (name) => {
	switch (name) {
		case "budget":
			return <Preferences id={User.GetId()} />;
		default:
			return null;
	}
};

const Plugins = ({
	error,
	loading,
	plugins,
	searchTerm,
	updateSearchTerm,
	handleSearch,
	handleEnablePlugin,
	handleDisablePlugin,
}) => {
	const PluginItem = ({
		displayName,
		name,
		id,
		description,
		consumedBy,
		consumedByUser,
	}) => (
		<>
			<Item>
				<Item.Content>
					<Item.Header as="h4">{displayName}</Item.Header>
					<Item.Meta>Description</Item.Meta>
					<Item.Description>{description}</Item.Description>
					<Item.Extra>
						<Label>
							{consumedBy} <Icon name="users" />
						</Label>

						<ButtonGroup size="small" floated="right">
							{consumedByUser ? (
								<>
									<Suspense>{getPluginModal(name)}</Suspense>
									<Button
										type="button"
										color="grey"
										onClick={() => handleDisablePlugin(id)}
									>
										<Icon name="minus" />
										Disable Plugin
									</Button>
								</>
							) : (
								<Button
									type="button"
									color="blue"
									onClick={() => handleEnablePlugin(id)}
								>
									<Icon name="plus" />
									Enable Plugin
								</Button>
							)}
						</ButtonGroup>
					</Item.Extra>
				</Item.Content>
			</Item>
		</>
	);

	const enabledPlugins = plugins.filter((x) => x.consumedByUser);
	const disabledPlugins = plugins.filter((x) => !x.consumedByUser);

	return (
		<>
			{error ? (
				<Message error header="An error occured!" content={error} />
			) : null}

			<Divider horizontal>
				<h3>My Plugins</h3>
			</Divider>

			{enabledPlugins.length > 0 ? (
				<Item.Group divided>
					{enabledPlugins.map((p, key) => (
						<PluginItem key={key} {...p} />
					))}
				</Item.Group>
			) : (
				<p>You haven't enabled any plugins, add some below!</p>
			)}

			{disabledPlugins.length > 0 ? (
				<>
					<Divider horizontal>
						<h3>More</h3>
					</Divider>

					<Form onSubmit={handleSearch}>
						<Form.Group inline>
							<Form.Field>
								<input
									placeholder="Search..."
									value={searchTerm}
									onChange={updateSearchTerm}
									type="text"
								/>
							</Form.Field>
							<Form.Field>
								<Button
									color="blue"
									type="submit"
									loading={loading}
								>
									Search
								</Button>
							</Form.Field>
						</Form.Group>
					</Form>

					<Item.Group>
						{disabledPlugins.map((p, key) => (
							<PluginItem key={key} {...p} />
						))}
					</Item.Group>
				</>
			) : null}
		</>
	);
};

export default Plugins;
