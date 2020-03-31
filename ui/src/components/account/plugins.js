import React from "react";
import {
	Item,
	Button,
	Icon,
	Form,
	Message,
	Divider,
	Grid,
} from "semantic-ui-react";

const Plugins = ({
	error,
	loading,
	plugins,
	searchTerm,
	updateSearchTerm,
	showMore,
	toggleMore,
	handleSearch,
	handleEnablePlugin,
	handleDisablePlugin,
}) => {
	const PluginItem = ({
		displayName,
		id,
		description,
		consumedBy,
		consumedByUser,
	}) => (
		<Item>
			<Item.Content>
				<Item.Header as="h4">{displayName}</Item.Header>
				<Item.Meta>Description</Item.Meta>
				<Item.Description>{description}</Item.Description>
				<Item.Extra>
					{consumedBy > 0
						? `This plugin is being used by ${consumedBy} other users.`
						: "There are currently no users using this plugin."}

					{consumedByUser ? (
						<Button
							type="button"
							color="grey"
							size="small"
							floated="right"
							onClick={() => handleDisablePlugin(id)}
						>
							<Icon name="minus" />
							Disable Plugin
						</Button>
					) : (
						<Button
							type="button"
							size="small"
							floated="right"
							color="blue"
							onClick={() => handleEnablePlugin(id)}
						>
							<Icon name="plus" />
							Enable Plugin
						</Button>
					)}
				</Item.Extra>
			</Item.Content>
		</Item>
	);

	const enabledPlugins = plugins.filter(x => x.consumedByUser);
	const disabledPlugins = plugins.filter(x => !x.consumedByUser);

	return (
		<>
			{error ? (
				<Message error header="An error occured!" content={error} />
			) : null}

			<Divider horizontal>
				<h3>Enabled</h3>
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
