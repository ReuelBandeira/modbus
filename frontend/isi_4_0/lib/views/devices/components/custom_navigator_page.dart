import 'package:flutter/material.dart';
import 'package:isi_4_0/colors/color.dart';
import 'package:isi_4_0/providers/lang.dart';

class CustomNavigatorPage extends StatelessWidget {
  const CustomNavigatorPage({
    super.key,
    required this.language,
    this.paddingRight,
    required this.page,
    required this.items,
    required this.totalItems,
    required this.totalPages,
    required this.onItemsChange,
    required this.onPageChange,
  });
  final LanguageProvider language;
  final double? paddingRight;
  final int page;
  final int items;
  final int totalItems;
  final int totalPages;
  final Function onItemsChange;
  final Function onPageChange;

  @override
  Widget build(BuildContext context) {
    Map<String, dynamic> lang =
        language.getDataLanguage(language.currentLanguage)['headDeviceList'];
    final translator = language.getDataLanguage(language.currentLanguage);
    String selectedValue = items.toString();
    List<String> lines = [
      '1',
      '2',
      '3',
      '4',
      '5',
      '10',
      '20',
      '30',
      '40',
      '50'
    ];

    return Container(
      height: 40,
      width: 180 * lang.length.toDouble(),
      alignment: Alignment.centerRight,
      padding: EdgeInsets.only(right: paddingRight ?? 0.0),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.end,
          children: [
            Text(
              translator['dictionary']['page_lines'],
              style: const TextStyle(fontSize: 12),
            ),
            Container(
              height: 40,
              width: 60,
              padding: const EdgeInsets.only(right: 15),
              child: DropdownButton<String>(
                value: selectedValue,
                isExpanded: true,
                underline: const SizedBox.shrink(),
                onChanged: (newValue) {
                  if (newValue != null) {
                    onItemsChange(int.parse(newValue));
                  } else {
                    onItemsChange(10);
                  }
                },
                items: lines.map<DropdownMenuItem<String>>((String value) {
                  return DropdownMenuItem<String>(
                    value: value,
                    child: Text(value, style: const TextStyle(fontSize: 12)),
                  );
                }).toList(),
              ),
            ),
            Container(
                padding: const EdgeInsets.only(right: 8.0),
                child: Text(
                  '${1 + ((page - 1) * items)}'
                      '-'
                      '${(page * items) < totalItems ? (page * items) : totalItems} '
                      '${translator['dictionary']['of']} $totalItems',
                  style: const TextStyle(
                      fontSize: 12, color: CustomColors.primaryColorApp),
                )),
            Align(
              alignment: Alignment.center,
              child: SizedBox(
                height: 35,
                width: 35,
                child: TextButton(
                  onPressed: (page == 1) ? null : () {
                    if (page != 1) {
                      onPageChange(1);
                    }
                  },
                  style: ButtonStyle(
                      padding:
                          MaterialStateProperty.all(const EdgeInsets.all(0)),
                      shape: MaterialStateProperty.all(RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(90)))),
                  child: const Icon(
                    Icons.first_page_outlined,
                    size: 18,
                  ),
                ),
              ),
            ),
            Align(
              alignment: Alignment.center,
              child: SizedBox(
                height: 35,
                width: 35,
                child: TextButton(
                  onPressed: (page == 1) ? null : () {
                    if (page > 1) {
                      onPageChange(page - 1);
                    }
                  },
                  style: ButtonStyle(
                      padding:
                          MaterialStateProperty.all(const EdgeInsets.all(0)),
                      shape: MaterialStateProperty.all(RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(90)))),
                  child: const Icon(
                    Icons.arrow_back_ios_outlined,
                    size: 11,
                  ),
                ),
              ),
            ),
            Align(
              alignment: Alignment.center,
              child: SizedBox(
                height: 35,
                width: 35,
                child: TextButton(
                  onPressed: ((page + 1) > totalPages) ? null : () {
                    if ((page + 1) <= totalPages) {
                      onPageChange(page + 1);
                    }
                  },
                  style: ButtonStyle(
                      padding:
                          MaterialStateProperty.all(const EdgeInsets.all(0)),
                      shape: MaterialStateProperty.all(RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(90)))),
                  child: const Icon(
                    Icons.arrow_forward_ios_outlined,
                    size: 11,
                  ),
                ),
              ),
            ),
            Align(
              alignment: Alignment.center,
              child: SizedBox(
                height: 35,
                width: 35,
                child: TextButton(
                  onPressed: (page == totalPages) ? null : () {
                    if (page != totalPages) {
                      onPageChange(totalPages);
                    }
                  },
                  style: ButtonStyle(
                      padding:
                          MaterialStateProperty.all(const EdgeInsets.all(0)),
                      shape: MaterialStateProperty.all(RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(90)))),
                  child: const Icon(
                    Icons.last_page_outlined,
                    size: 18,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
