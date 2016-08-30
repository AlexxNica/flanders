//Adding the string for +1 into text field by button click.
$(function () {
  $('#button').on('click', function () {
    var text = $('#text');
    text.val(text.val() + '%2b1');
  });
});
