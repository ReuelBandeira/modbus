class DataDevices {
  String id;
  List<String> fields;
  List<String> topics;
  List<dynamic> extraFields;
  List<dynamic> extraFieldsData;

  DataDevices(
      {required this.id,
      required this.fields,
      required this.topics,
      this.extraFields = const [],
      this.extraFieldsData = const []});
}
