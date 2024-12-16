import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "explorePage.dart";
import '../globalData.dart';

class TopBarExplore extends StatelessWidget {
  TopBarExplore({
    super.key,
    required this.currentIndex,
  });

  final int currentIndex;
  final List<Widget> _pages = [
    ExplorePage(explorePageIndex: 0,),
    ExplorePage(explorePageIndex: 1,),
    ExplorePage(explorePageIndex: 2,),
  ];

  Future<void> changePage(int index, BuildContext context) async {
    if (index != currentIndex) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        Navigator.push(
          Globaldata.myContext,
          MaterialPageRoute(builder: (context) => _pages[index]),
        );
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        TextButton(
          onPressed: () {
            changePage(0, Globaldata.myContext);
          },
          child: AutoSizeText(
            "All",
            style: TextStyle(
                fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black),
            maxLines: 1,
          ),
        ),
        TextButton(
          onPressed: () {
            changePage(1, Globaldata.myContext);
          },
          child: AutoSizeText(
            "Workflows",
            style: TextStyle(
                fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black),
            maxLines: 1,
          ),
        ),
        TextButton(
          onPressed: () {
            changePage(2, Globaldata.myContext);
          },
          child: AutoSizeText(
            "Services",
            style: TextStyle(
                fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black),
            maxLines: 1,
          ),
        ),
      ],
    );
  }
}
