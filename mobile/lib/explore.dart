import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:math';

import "custom_widget.dart";
import "footer_page.dart";

class ExplorePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView (
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            ExploreTopBar(),
            MyDivider(),
            SearchBar(),
            MyCard(),
            Padding(padding: EdgeInsets.only(top: 10)),
            MyCard(),
            Padding(padding: EdgeInsets.only(top: 10)),
            MyCard(),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 1),
    );
  }
}

class MyCard extends StatelessWidget {
  const MyCard({super.key});

  @override
  Widget build(BuildContext context) {
    List<Color> colorList = [
      const Color.fromARGB(255, 46, 41, 78),
      const Color.fromARGB(255, 44, 110, 252),
      const Color.fromARGB(255, 228, 64, 95),
      const Color.fromARGB(255, 29, 185, 84),
      const Color.fromARGB(255, 205, 65, 43),
    ];

    int colorIndex = Random().nextInt(colorList.length);

    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return Container(
            decoration: BoxDecoration(
              border: Border.all(
                color: colorList[colorIndex],
                width: 2.0,
              ),
              borderRadius: BorderRadius.circular(4.0),
            ),
            child: InkWell(
              onTap: () {
                debugPrint('Card tapped.');
              },
              child: Column(
                children: [
                  SizedBox(
                    width: constraints.maxWidth,
                    height: 20,
                    child: DecoratedBox(
                      decoration: BoxDecoration(
                        color: colorList[colorIndex],
                      ),
                    ),
                  ),
                  Padding(
                    padding: EdgeInsets.symmetric(horizontal: 5, vertical: 10),
                    child: SizedBox(
                    width: constraints.maxWidth,
                    child: Text( "A card that can be tapped\nAnd how fit with his content",
                        style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.black,),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class SearchBar extends StatelessWidget {
  const SearchBar({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: Container(
        decoration: BoxDecoration(
          color: const Color.fromARGB(255, 238, 238, 238),
          borderRadius:  BorderRadius.circular(5),
        ),
        child: TextField(
          decoration: InputDecoration(
            prefixIcon: Padding(
              padding: EdgeInsets.only(left: 5),
              child: Icon(
                Icons.search,
                size: 35,
                color: Colors.black,
              ),
            ),
            hintStyle: TextStyle(
              fontSize: 20,
              color: const Color.fromARGB(255, 156, 156, 156),
              fontWeight: FontWeight.bold,
            ),
            hintText: 'Search',
            border: InputBorder.none,
            contentPadding: EdgeInsets.symmetric(vertical: 10, horizontal: 10),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(10),
              borderSide: BorderSide(
                color: Colors.black,
                width: 4,
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class ExploreTopBar extends StatelessWidget {
  const ExploreTopBar({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Row (
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        TextButton(
          onPressed: () {},
          child: AutoSizeText("All", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),

        TextButton(
          onPressed: () {},
          child: AutoSizeText("Workflow", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),

        TextButton(
          onPressed: () {},
          child: AutoSizeText("Services", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),

        TextButton(
          onPressed: () {},
          child: AutoSizeText("Stories", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),
      ],
    );
  }
}
