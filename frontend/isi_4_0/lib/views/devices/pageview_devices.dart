import 'package:flutter/material.dart';
import 'package:isi_4_0/views/devices/devices.dart';
import 'package:isi_4_0/views/devices/pageview_extra_data.dart';

class PageViewDevices extends StatefulWidget {
  const PageViewDevices({super.key});
  @override
  State<PageViewDevices> createState() => _PageViewDevicesState();
}

class _PageViewDevicesState extends State<PageViewDevices> {
  PageController pageController = PageController();

  @override
  Widget build(BuildContext context) {
    return PageView(
      physics: const NeverScrollableScrollPhysics(),
      controller: pageController,
      scrollDirection: Axis.horizontal,
      children: [
        Devices(
          onTap: () {
            pageviewNavigator(pageController, 1);
          },
        ),
        PageViewExtraData(
          onTap: () {
            pageviewNavigator(pageController, 0);
          },
        )
      ],
    );
  }

  void pageviewNavigator(PageController pageController, int position) {
    pageController.animateToPage(position,
        duration: const Duration(milliseconds: 1000), curve: Curves.easeIn);
  }
}
