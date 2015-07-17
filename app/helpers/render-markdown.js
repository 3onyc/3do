import Ember from 'ember';

export default Ember.Handlebars.makeBoundHelper( function(value) {
    var writer = new Remarkable("commonmark");
    var renderedHtml = writer.render(value);

    return new Ember.Handlebars.SafeString(renderedHtml);
});
