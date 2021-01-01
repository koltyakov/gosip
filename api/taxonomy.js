var initTaxonomy = (callback) => {
  const getScript = (source, callback) => {
    const script = document.createElement('script');
    const prior = document.getElementsByTagName('script')[0];
    script.async = 1;
    script.onload = script.onreadystatechange = (_, isAbort) => {
      if(isAbort || !script.readyState || /loaded|complete/.test(script.readyState)) {
        script.onload = script.onreadystatechange = null;
        script = undefined;
        if(!isAbort && callback) setTimeout(callback, 0);
      }
    }
    script.src = source;
    prior.parentNode.insertBefore(script, prior);
  };

  const scriptbase = `${_spPageContextInfo.siteAbsoluteUrl}/_layouts/15`;
  if(typeof SP.Taxonomy === 'undefined') {
    getScript(`${scriptbase}/SP.Taxonomy.js`, callback);
  } else {
    callback();
  }
};

var executeQueryAsync = (context) => {
  return new Promise((resolve, reject) => {
    context.executeQueryAsync(resolve, (_, args) => {
      reject(args.get_message());
    });
  });
};

// Modify with CSOM object calls, run and check packages in browser dev tools
initTaxonomy(async () => {
  const ctx = SP.ClientContext.get_current();
  const taxSession = SP.Taxonomy.TaxonomySession.getTaxonomySession(ctx);
  const termStore = taxSession.getDefaultSiteCollectionTermStore();

  ctx.load(termStore);

  await executeQueryAsync(ctx);
  console.log(termStore);
});