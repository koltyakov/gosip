var executeQueryAsync = (context) => {
  return new Promise((resolve, reject) => {
    context.executeQueryAsync(resolve, (_, args) => {
      reject(args.get_message());
    });
  });
};

SP.SOD.executeFunc('sp.js', 'SP.ClientContext', () => {
  SP.SOD.registerSod('sp.taxonomy.js', SP.Utilities.Utility.getLayoutsPageUrl('sp.taxonomy.js'));
  SP.SOD.executeFunc('sp.taxonomy.js', 'SP.Taxonomy.TaxonomySession', async () => {

    const ctx = SP.ClientContext.get_current();
    const taxSession = SP.Taxonomy.TaxonomySession.getTaxonomySession(ctx);

    const termStore = taxSession.getDefaultSiteCollectionTermStore();

    ctx.load(termStore);

    await executeQueryAsync(ctx); console.clear();

    console.log(termStore);

  });
});